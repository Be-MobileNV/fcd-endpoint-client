import {ClientOptions, FcdWebSocketClient, GPSPosition} from "../../client/TypeScript/client";

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

runClient();
