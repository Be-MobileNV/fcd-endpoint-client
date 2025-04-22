import WebSocket from 'ws';
import { Buffer } from 'buffer';

interface GPSPosition {
  vehicleId: string;
  vehicleType?: number;
  engineState?: number;   // -1, 0, or 1
  timestamp: number;      // Unix milliseconds
  lon: number;
  lat: number;
  heading?: number;       // degrees
  hdop?: number;          // meters
  speed?: number;         // km/h
  alt?: number;           // meters
  metadata?: Record<string, string>;
}

interface ClientOptions {
  address: string;
  port: string;
  tls?: boolean;
  username?: string;
  password?: string;
  compression?: boolean;
}

class FcdWebSocketClient {
  private ws: WebSocket | null = null;
  private readonly url: string;
  private readonly options: ClientOptions;
  private readonly headers: Record<string, string> = {};
  private readonly connectionPromise: Promise<void>;
  private resolveConnection!: () => void;
  private rejectConnection!: (reason?: any) => void;


  constructor(options: ClientOptions) {
    this.options = options;
    const scheme = options.tls ? 'wss' : 'ws';
    this.url = `${scheme}://${options.address}:${options.port}/v1/ws`;

    if (options.username && options.password) {
      const credentials = Buffer.from(`${options.username}:${options.password}`).toString('base64');
      this.headers['Authorization'] = `Basic ${credentials}`;
    }

    this.connectionPromise = new Promise((resolve, reject) => {
        this.resolveConnection = resolve;
        this.rejectConnection = reject;
    });
  }

  connect(): Promise<void> {
    if (this.ws && (this.ws.readyState === WebSocket.OPEN || this.ws.readyState === WebSocket.CONNECTING)) {
      return this.connectionPromise;
    }

    const wsOptions: WebSocket.ClientOptions = {
      headers: this.headers,
      perMessageDeflate: this.options.compression ?? false,
    };

    console.log(`Connecting to ${this.url}...`);
    this.ws = new WebSocket(this.url, wsOptions);

    this.ws.on('open', () => {
      console.log('WebSocket connection established.');
      this.resolveConnection();
    });

    this.ws.on('message', (data) => {
      console.error(`Received error from server: ${data.toString()}`);
    });

    this.ws.on('error', (error) => {
      console.error('WebSocket error:', error);
      if (!this.connectionPromise) { // Reject initial connection promise if error occurs before open
        this.rejectConnection(error);
      }
      this.ws = null;
    });

    this.ws.on('close', (code, reason) => {
      console.log(`WebSocket connection closed. Code: ${code}, Reason: ${reason.toString()}`);
      this.ws = null;
    });

    return this.connectionPromise;
  }

  async sendPosition(position: GPSPosition): Promise<void> {
    // Basic validation could be added here if desired, similar to Go's Validate
    await this.connectionPromise;

    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.error('WebSocket is not open. Cannot send message.');
      throw new Error('WebSocket not open');
    }

    const payload = JSON.stringify(position);
    this.ws.send(payload, (error) => {
      if (error) {
        console.error('WebSocket send error:', error);
      }
    });
  }

  close(code: number = 1000, reason: string = 'Client closing'): void {
    if (this.ws) {
      console.log(`Closing WebSocket connection with code ${code}...`);
      this.ws.close(code, reason);
      this.ws = null;
    }
  }
}


// --- Example Usage ---
async function runClient() {
  const config: ClientOptions = {
    address: '127.0.0.1',
    port: '8080',
    tls: false,
    username: '', // Add credentials if needed
    password: '',
    compression: true,
  };

  const client = new FcdWebSocketClient(config);

  try {
    await client.connect();

    // Send some sample positions
    for (let i = 0; i < 100; i++) {
      const position: GPSPosition = {
        vehicleId: `TS-VEHICLE-${Math.floor(Math.random() * 1000)}`,
        timestamp: Date.now(),
        lon: 4.35 + (Math.random() - 0.5) * 0.1,
        lat: 50.85 + (Math.random() - 0.5) * 0.1,
        heading: Math.random() * 360,
        speed: Math.random() * 100,
      };
      await client.sendPosition(position);
      await new Promise(resolve => setTimeout(resolve, 250));
    }

  } catch (error) {
    console.error("Client operation failed:", error);
  } finally {
    // Close the connection after sending
    client.close();
    console.log("Client finished.");
  }
}

// Run the example client
runClient();