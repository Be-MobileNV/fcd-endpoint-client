import argparse
import json

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