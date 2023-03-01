import base64
import logging
import json
import websockets
import sys
import asyncio


class GPSPosition:
    def __init__(self, vehicle_id, timestamp, lon, lat, vehicle_type=None, heading=None, hdop=None, speed=None, engine_state=None):
        self.vehicle_id = str(vehicle_id)
        self.timestamp = int(timestamp.timestamp() * 1000)  # Should be a datetime.datetime instance in UTC
        self.lon = float(lon)
        self.lat = float(lat)
        self.vehicle_type = int(vehicle_type) if vehicle_type is not None else None
        self.heading = float(heading) if heading is not None else None
        self.hdop = float(hdop) if hdop is not None else None
        self.speed = float(speed) if speed is not None else None
        self.engine_state = int(engine_state) if engine_state is not None else None

    def validate(self):
        if not self.vehicle_id:
            raise ValueError("empty vehicle ID")
        if not len(self.vehicle_id) <= 64:
            raise ValueError(f"vehicle ID too long: {self.vehicle_id}")
        if not -180 <= self.lon <= 180:
            raise ValueError(f"invalid lon: {self.lon}")
        if not -90 <= self.lat <= 90:
            raise ValueError(f"invalid lat: {self.lat}")
        if not (self.vehicle_type is None or self.vehicle_type >= 0):
            raise ValueError(f"invalid vehicle type: {self.vehicle_type}")
        if not (self.heading is None or 0 <= self.heading < 360):
            raise ValueError(f"invalid heading: {self.heading}")
        if not (self.hdop is None or self.hdop >= 0):
            raise ValueError(f"invalid hdop: {self.hdop}")
        if not (self.speed is None or self.speed >= 0):
            raise ValueError(f"invalid speed: {self.speed}")
        if not (self.engine_state is None or -1 <= self.engine_state <= 1):
            raise ValueError(f"invalid engine state: {self.engine_state}")

    def to_json(self):
        # Mandatory fields
        pos = {"vehicleId": self.vehicle_id,
               "timestamp": self.timestamp,
               "lon": self.lon,
               "lat": self.lat}

        # Optional fields
        if self.vehicle_type is not None:
            pos["vehicleType"] = self.vehicle_type
        if self.heading is not None:
            pos["heading"] = self.heading
        if self.hdop is not None:
            pos["hdop"] = self.hdop
        if self.speed is not None:
            pos["speed"] = self.speed
        if self.engine_state is not None:
            pos["engineState"] = self.engine_state

        return json.dumps(pos)


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
            self.logger.info("Opened websocket connection")
            task = asyncio.create_task(self.handle_replies(websocket))
            async for pos in generator:
                try:
                    pos.validate()
                except Exception as E:
                    self.logger.error(E)
                    continue

                self.logger.info("Sending GPS position")
                # We need to wrap websocket.send(pos) in a Task in order to handle exceptions
                sendtask = asyncio.create_task(websocket.send(pos.to_json()))
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
