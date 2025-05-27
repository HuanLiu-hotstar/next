const { ApolloServer, gql } = require('apollo-server');
const { buildFederatedSchema } = require('@apollo/federation');

// Define schema
const typeDefs = gql `
  type Query {
    fetchContentCollection(tray_id: String!): ContentCollection
  }

  type CollectionItems @key(fields: "content_id") {
    content_id: String!
    name: String
    title: String
	channel: Channel @requires(fields: "content_id")
  }

   # This type is provided by the Go subgraph
   extend type Channel @key(fields: "content_id") {
    content_id: String! @external
  }


  type ContentCollection {
    tray_id: String!
    tray_name: String
    collectionItems: [CollectionItems!]
  }
`;

// Sample data
const collections = [{
    tray_id: 'tray123',
    tray_name: 'Popular Shows',
    collectionItems: [
        { content_id: 'content1', name: 'Item A', title: 'Title A' },
        { content_id: 'content2', name: 'Item B', title: 'Title B' },
        { content_id: 'content3', name: 'Item C', title: 'Title C' },
        { content_id: 'content4', name: 'Item D', title: 'Title D' },
        { content_id: 'content6', name: 'Item F', title: 'Title F' },
    ]
}, {
    tray_id: 'tray456',
    tray_name: 'Popular Movie',
    collectionItems: [
        { content_id: 'content4', name: 'Item D', title: 'Title D' },
        { content_id: 'content6', name: 'Item F', title: 'Title F' },
    ]
}];

// Define resolvers
const resolvers = {
    Query: {
        fetchContentCollection: (_, { tray_id }) => {
            return collections.find(c => c.tray_id === tray_id);
        }
    },
    CollectionItems: {
        __resolveReference(ref) {
            return collections
                .flatMap(c => c.collectionItems)
                .find(item => item.content_id === ref.content_id);
        },
        // This is a placeholder; actual `channel` resolution is in the Go subgraph
        channel(parent) {
            // Do nothing – let the federated gateway handle it
            return { content_id: parent.content_id };
        }
    }
};

// Create and start server
const server = new ApolloServer({
    schema: buildFederatedSchema([{ typeDefs, resolvers }])
});

server.listen({ port: 4002 }).then(({ url }) => {
    console.log(`🚀 Content subgraph running at ${url}`);
});