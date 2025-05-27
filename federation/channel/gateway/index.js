const { ApolloServer } = require('@apollo/server');
const { startStandaloneServer } = require('@apollo/server/standalone');
const { ApolloGateway, IntrospectAndCompose } = require('@apollo/gateway');

const gateway = new ApolloGateway({
    supergraphSdl: new IntrospectAndCompose({
        subgraphs: [
            { name: 'channels', url: 'http://localhost:4001/query' },
            { name: 'content', url: 'http://localhost:4002/query' },
        ],
    }),
});

const server = new ApolloServer({
    gateway,
    subscriptions: false,
});

// Use an async function instead of top-level await
async function startServer() {
    try {
        const { url } = await startStandaloneServer(server);
        console.log(`🚀  Server ready at ${url}`);
    } catch (error) {
        console.error('Error starting server:', error);
    }
}

// Start the server
startServer();