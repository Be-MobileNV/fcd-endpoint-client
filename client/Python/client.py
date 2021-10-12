import time
import base64
import logging
import argparse
import json
import websockets

close_wait_time = 10

class gpsPosition:
    def __init__(self, vehicleID, vehicleType, timestamp, lon, lat, heading, hdop, speed):
        self.vehicleID = vehicleID
        self.vehicleType = vehicleType
        self.timestamp = timestamp
        self.lon = lon
        self.lat = lat
        self.heading = heading
        self.hdop = hdop
        self.speed = speed

    def toJson(self):
        gps_position = {"vehicleId": self.vehicleID,
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
        self.logger = self.setLogger(self, args)
        self.url = getURL(args)
        self.username = args.username
        self.password = args.password
        self.client = self.connect()
    
    def setLogger(self, args):
        logger = logging.getLogger("FCD Endpoint Websocket Client")
        if args.verbose is True:
            logging.basicConfig(level=logging.DEBUG)
            logger.setLevel(logging.DEBUG)
        else:
            logging.basicConfig(level=logging.INFO)
        return logger

    def connect(self):
        header = []
        if self.username != "" and self.password != "":
            basic_auth_str = self.username + ":" + self.password
            header = ["Authorization: Basic {}".format(base64.b64encode(basic_auth_str.encode('ascii')).decode())]
        self.logger.info("Connecting to {}".format(self.url))
        return websockets.client.connect(self.url, extra_headers=header, ping_interval=60, ping_timeout=10, close_timeout=10)
    
    def sendGPSPosition(self, position):
        self.logger.debug("GPS position JSON: {}".format(position))
        self.logger.info("Sending GPS position")
        self.client.send(position)
        self.logger.info("Waiting {} seconds on server to receive possible errors".format(close_wait_time))
        time.sleep(close_wait_time)  # Sleep is mandatory, otherwise errors are not received anymore from the server
    
    def close(self):
        self.logger.info("Closing the websocket by sending a close message")
        self.client.close()

def loadConfig():
    parser = argparse.ArgumentParser(description="FCD Endpoint Websocket Client")
    parser.add_argument('--address', type=str, default='127.0.0.1',
                        help='The address of the server to send the data to')
    parser.add_argument('--port', type=str, default='443',
                        help='If ingress is between the client and server, use 443, '
                            'otherwise the same port as the server')
    parser.add_argument('--tls', type=str, default="true",
                        help='Usage of WSS(true) or WS(false).')
    parser.add_argument('--username', type=str, default='',
                        help='The username if you want to use basic authorization (securews must be set on true)')
    parser.add_argument('--password', type=str, default='',
                        help='The password if you want to use basic authorization (securews must be set on true)')

    parser.add_argument('--verbose', default=False, action="store_true", help="Increase log verbosity", )
    return parser.parse_args()

def initiateWebSocket():
    args = loadConfig()
    ws = WebSocket(args)
    return ws

def getURL(args):
    if args.tls in ["true", "True"]:
        # Secured WSS protocol
        url = "wss://{}:{}/v1/ws".format(args.address, args.port)
    else:
        # Unsecured WS protocol
        url = "ws://{}:{}/v1/ws".format(args.address, args.port)
    return url