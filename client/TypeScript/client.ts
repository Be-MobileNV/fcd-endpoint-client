import WebSocket from 'ws';
import { Buffer } from 'buffer';

export interface GPSPosition {
  vehicleId: string;
  vehicleType?: number;
  engineState?: number;   // -1, 0, or 1
  timestamp: number;      // Unix milliseconds
  lon: number;
  lat: number;
  heading?: number;       // degrees
  hacc?: number;          // meters
  speed?: number;         // km/h
  alt?: number;           // meters
  metadata?: Record<string, string>;
}

export interface ClientOptions {
  address: string;
  port: string;
  tls?: boolean;
  username?: string;
  password?: string;
  compression?: boolean;
}

export class FcdWebSocketClient {
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
