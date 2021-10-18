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

namespace WSClient
{
    public class GPSPosition {
        [JsonPropertyName("vehicleId")]
        public string VehicleID { get; set; }
        [JsonPropertyName("vehicleType")]
        public int VehicleType { get; set; }
        [JsonPropertyName("engineState")]
        public int EngineState { get; set; }
        [JsonPropertyName("timestamp")]
        public Int64 Timestamp { get; set; }
        [JsonPropertyName("lon")]
        public double Longitude  { get; set; }
        [JsonPropertyName("lat")]
        public double Latitude { get; set; }
        [JsonPropertyName("heading")]
        public float Heading { get; set; }
        [JsonPropertyName("hdop")]
        public float HDOP { get; set; }
        [JsonPropertyName("speed")]
        public float Speed { get; set; }
        [JsonPropertyName("metadata")]
        public Dictionary<string, string> Metadata { get; set; }
    }
    public class Client
    {
        public string Address { get; } = "";
        public int Port { get; } = 443;
        public bool TLS { get; } = true;
        public string Username { get; } = "";
        public string Password { get; } = "";
        public string url;
        public ClientWebSocket ws; 

        public Client(string Address, int Port, string Username, string Password, bool TLS){
            this.Address = Address;
            this.Port = Port;
            this.Username = Username;
            this.Password = Password;
            this.TLS = TLS;
            this.url = getURL();
        }

        public void send(GPSPosition position){
            JsonSerializerOptions options = new JsonSerializerOptions()
            {
                DefaultIgnoreCondition = JsonIgnoreCondition.WhenWritingNull
            };

            string message = JsonSerializer.Serialize<GPSPosition>(position, options);
            var bytesToSend = new ArraySegment<byte>(Encoding.UTF8.GetBytes(message));
            Task.Run(async () =>
            {
                try
                {
                    using (var socket = new ClientWebSocket()) {
                        socket.Options.KeepAliveInterval = TimeSpan.FromSeconds(5); // keep alive interval (ping pong) - https://github.com/dotnet/runtime/blob/7cbf0a7011813cb84c6c858ef19acb770daa777e/src/libraries/Common/src/System/Net/WebSockets/ManagedWebSocket.cs#L886
                        socket.Options.Credentials = new NetworkCredential(Username, Password);
                        await socket.ConnectAsync(new Uri(url), CancellationToken.None);
                        var tSend = sendPoint(socket, bytesToSend);

                        await socket.CloseOutputAsync(WebSocketCloseStatus.NormalClosure, "", CancellationToken.None);
                        Console.WriteLine("Done!");
                    }
                }
                catch (System.Exception ex)
                {
                    Console.WriteLine($"ERROR setting up socket - {ex.Message}");
                    return;
                }

            }).GetAwaiter().GetResult();
        }

        public static async Task sendPoint(ClientWebSocket socket, ArraySegment<byte> bytes)
        {
            try
            {
                Console.WriteLine($"Sending GPS position");
                await socket.SendAsync(bytes, WebSocketMessageType.Text, true, CancellationToken.None);
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Failed to send GPS position message to server: {ex.Message}");
                return;
            }
        }

        private string getURL(){
            if (TLS){
                return string.Format("wss://{0}:{1}/v1/ws", Address, Port);
            } else {
                return string.Format("ws://{0}:{1}/v1/ws", Address, Port);
            }
        }
    }
}
