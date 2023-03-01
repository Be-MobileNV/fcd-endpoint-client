import sys
import os
import datetime
import random
import string
import asyncio
import argparse


sys.path.append(os.path.join(os.path.dirname(__file__), '..', '..', 'client', 'Python'))
import client


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


async def main():
    args = loadConfig()
    ws = client.WebSocket(args)
    await ws.send(generate_positions())


async def generate_positions(n=100):
    """
    Generate n random positions
    """
    for i in range(n):
        yield get_position()


def get_position():
    ymin, xmin = 46.691265, 4.565761
    ymax, xmax = 52.076458, 6.257655
    pos = client.GPSPosition(
        vehicle_id=get_vehicle_id(),
        timestamp=datetime.datetime.now(datetime.timezone.utc),
        lon=random.uniform(xmin, xmax),
        lat=random.uniform(ymin,ymax),
        heading=random.uniform(0, 360),
        hdop=random.uniform(0, 10),
        speed=None, # In this example, the speed is unknown, so we set it to None
        vehicle_type=1,
    )
    return pos


def get_vehicle_id():
    charset = string.ascii_letters + string.digits + "-"
    vID = ''.join(random.choice(charset) for i in range(10))
    return vID


if __name__ == "__main__":
    asyncio.run(main())
