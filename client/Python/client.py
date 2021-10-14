import base64
import logging
import json
import websockets
import asyncio


class GPSPosition:
    def __init__(self, vehicle_id, vehicle_type, timestamp, lon, lat, heading, hdop, speed):
        self.vehicleId = vehicle_id
        self.vehicleType = vehicle_type
        self.timestamp = timestamp
        self.lon = lon
        self.lat = lat
        self.heading = heading
        self.hdop = hdop
        self.speed = speed

    def toJSON(self):
        gps_position = {"vehicleId": self.vehicleId,
                        "vehicleType": self.vehicleType,
                        "timestamp": self.timestamp,
                        "lon": self.lon,
                        "lat": self.lat,
                        "heading": self.heading,
                        "hdop": self.hdop,
                        "speed": self.speed,
                        }
        return json.dumps(gps_position)


class WebSocket:
    def __init__(self, args):
        self.logger = self._logger(args)
        self.url = self._url(args)
        self.username = args.username
        self.password = args.password

    @staticmethod
    def _url(args):
        if args.tls in ["true", "True"]:
            # Secured WSS protocol
            url = "wss://{}:{}/v1/ws".format(args.address, args.port)
        else:
            # Unsecured WS protocol
            url = "ws://{}:{}/v1/ws".format(args.address, args.port)
        return url
    
    @staticmethod
    def _logger(args):
        logger = logging.getLogger("FCD Endpoint Websocket Client")
        if args.verbose is True:
            logging.basicConfig(level=logging.DEBUG)
            logger.setLevel(logging.DEBUG)
        else:
            logging.basicConfig(level=logging.INFO)
        logger.info("Loaded config")
        return logger

    async def send(self, generator):
        header = []
        if self.username != "" and self.password != "":
            basic_auth_str = self.username + ":" + self.password
            header = [("Authorization", "Basic {}".format(base64.b64encode(basic_auth_str.encode('ascii')).decode()))]
        self.logger.info("Connecting to {}".format(self.url))
        async for websocket in websockets.connect(self.url, extra_headers=header, ping_interval=60, ping_timeout=10, close_timeout=10):
            task = asyncio.create_task(self.handle_replies(websocket))
            async for pos in generator:
                self.logger.info("Sending GPS position")
                # We need to wrap websocket.send(pos) in a Task in order to handle exceptions
                sendtask = asyncio.create_task(websocket.send(pos))
                try:
                    await sendtask
                except Exception as E:
                    self.logger.error(E)
                    break  # Force a reconnection
            else:
                # generator finished, exit
                break
            task.cancel()

    async def handle_replies(self, websocket):
        """
        Handle replies from the FCD endpoint server by logging them. Normally,
        there are no replies: the only replies are error messages.
        """
        async for reply in websocket:
            self.logger.error(reply)
