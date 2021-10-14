import sys, os, time, random, string, asyncio

sys.path.append(os.path.join(os.path.dirname(__file__), '..', '..', 'client', 'Python'))
import client

def main():
    ws = client.initiateWebSocket()
    print(ws)
    time.sleep(0.25)


def getPosition():
    ymin,xmin = 46.691265,4.565761
    ymax,xmax = 52.076458,6.257655
    pos = client.gpsPosition(getVehicleID(), 1, int(time.time() * 1000), random.uniform(xmin, xmax), random.uniform(ymin,ymax), random.uniform(0, 360), random.uniform(0,10), random.uniform(0, 120))
    return pos.toJson()

def getVehicleID():
    charset = string.ascii_letters + string.digits + "-"
    vID = ''.join(random.choice(charset) for i in range(10))
    return vID

if __name__ == "__main__":
    main()