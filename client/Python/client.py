import time
import base64
import logging
from client.Python.config import gpsPosition
import config
import websockets

try:
    import thread
except ImportError:
    import _thread as thread

close_wait_time = 10
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

def initiateWebSocket():
    args = config.loadConfig()
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