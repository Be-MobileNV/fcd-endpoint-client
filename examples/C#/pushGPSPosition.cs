using WSClient; 
using System;
using System.Collections.Generic;
using System.Text.Json;
using System.Text.Json.Serialization;
using System.Net.WebSockets;
using System.Net;
using System.Threading;
using System.Text;
using System.Threading.Tasks;
using McMaster.Extensions.CommandLineUtils;

namespace ClientTest
{
    public class Test{
        public static int Main(string[] args) => CommandLineApplication.Execute<Test>(args);

        [Option(ShortName= "a", LongName= "address", Description = "The address of the server to send the data to")]
        public string Address { get; } = "127.0.0.1";

        [Option(ShortName= "p", LongName= "port", Description = "If ingress is between the client and server, use 443, otherwise the same port as the server")]
        public int Port { get; } = 443;

        [Option(ShortName= "s", LongName= "securews", Description = "Usage of WSS(true) or WS(false).")]
        public bool TLS { get; } = true;
    
        [Option(ShortName= "u", LongName= "username", Description = "The username if you want to use basic authorization (securews must be set on true)")]
        public string Username { get; } = "";

        [Option(ShortName= "pwd", LongName= "password", Description = "The password if you want to use basic authorization (securews must be set on true)")]
        public string Password { get; } = "";

        public WSClient.Client client;

        private void OnExecute()
        {
            client = new WSClient.Client(Address, Port, Username, Password, TLS);

            for (int i = 0; i < 100; i++){
                client.send(getPosition());
            }
            
        }

        private WSClient.GPSPosition getPosition(){
            Random random = new Random();
            return new WSClient.GPSPosition(){
                VehicleID = getVehicleID(),
                VehicleType = 1,
                EngineState = 1,
                Timestamp = (long)(DateTime.UtcNow.Subtract(new DateTime(1970, 1, 1))).TotalMilliseconds,
                Longitude = random.NextDouble() * (6.257655-4.565761) + 4.565761,
                Latitude  = random.NextDouble() * (52.076458-46.691265) + 46.691265,
                Heading = (float)random.NextDouble() * 360,
                HDOP = (float)random.NextDouble(),
                Speed  = (float)random.NextDouble() * 120,
            };
        }

        private string getVehicleID(){
            var chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-";
            var stringChars = new char[12];
            var random = new Random();

            for (int i = 0; i < stringChars.Length; i++)
            {
                stringChars[i] = chars[random.Next(chars.Length)];
            }

            return new String(stringChars);
        }

    }
}