const express = require('express');
const { graphqlHTTP } = require('express-graphql');
const { buildSchema } = require('graphql');
const { buildSubgraphSchema } = require('@apollo/subgraph');

const fetch = require('node-fetch');

// Define schema with federation directives
const schema = buildSubgraphSchema(`
  scalar _Any
  scalar _FieldSet
  
  directive @external on FIELD_DEFINITION
  directive @requires(fields: _FieldSet!) on FIELD_DEFINITION
  directive @provides(fields: _FieldSet!) on FIELD_DEFINITION
  directive @key(fields: _FieldSet!) on OBJECT | INTERFACE
  
  type Content @key(fields: "content_id") {
    content_id: String!
    title: String
    name: String
    description: String
    tray_id: String
  }

  type _Service {
    sdl: String
  }

  type Query {
    hello: String
    greet(name: String!): String
    fetchContent(tray_id: String!): [Content]
    fetchContentV2(content_id: String!): Content
    _service: _Service
    _entities(representations: [_Any!]!): [_Entity]!
  }

  union _Entity = Content
`);

// Mock data for content
const mockContents = {
    "tray1": [{
            content_id: "content1",
            title: "First Content",
            name: "Content Series 1",
            description: "This is the first content in tray 1",
            tray_id: "tray1"
        },
        {
            content_id: "content2",
            title: "Second Content",
            name: "Content Series 1",
            description: "This is the second content in tray 1",
            tray_id: "tray1"
        },
        {
            content_id: "content3",
            title: "Third Content",
            name: "Content Series 2",
            description: "This is the third content in tray 1",
            tray_id: "tray1"
        },
        {
            content_id: "content4",
            title: "Fourth Content",
            name: "Content Series 2",
            description: "This is the fourth content in tray 1",
            tray_id: "tray1"
        }
    ],
    "tray2": [{
            content_id: "content5",
            title: "First Content in Tray 2",
            name: "Content Series 3",
            description: "This is the first content in tray 2",
            tray_id: "tray2"
        },
        {
            content_id: "content6",
            title: "Second Content in Tray 2",
            name: "Content Series 3",
            description: "This is the second content in tray 2",
            tray_id: "tray2"
        },
        {
            content_id: "content7",
            title: "Third Content in Tray 2",
            name: "Content Series 4",
            description: "This is the third content in tray 2",
            tray_id: "tray2"
        },
        {
            content_id: "content8",
            title: "Fourth Content in Tray 2",
            name: "Content Series 4",
            description: "This is the fourth content in tray 2",
            tray_id: "tray2"
        }
    ]
};

// Define resolvers
const root = {
    hello: () => 'Hello world!',
    greet: ({ name }) => `Hello, ${name}!`,
    fetchContent: ({ tray_id }) => {
        // Return mock data based on tray_id
        if (mockContents[tray_id]) {
            return mockContents[tray_id];
        }

        // If tray_id doesn't exist, return empty array
        return [];
    },
    fetchContentV2: ({ content_id }) => {
        // Find content by content_id across all trays
        for (const trayId in mockContents) {
            for (const content of mockContents[trayId]) {
                if (content.content_id === content_id) {
                    return content;
                }
            }
        }
        // Return null if content_id doesn't exist
        return null;
    },
    // Federation resolvers
    _service: () => ({
        sdl: schema.toString()
    }),
    _entities: ({ representations }) => {
        return representations.map(ref => {
            if (ref.__typename === 'Content') {
                const content_id = ref.content_id;
                // Find content by content_id
                for (const trayId in mockContents) {
                    for (const content of mockContents[trayId]) {
                        if (content.content_id === content_id) {
                            return content;
                        }
                    }
                }
            }
            // Filter out null values to avoid federation errors
            return null;
        }).filter(Boolean); // Filter out null values
    }
};

// Create an Express app
const app = express();

// Define the GraphQL endpoint
app.use('/query', graphqlHTTP({
    schema: schema,
    rootValue: root,
    graphiql: true, // Enable GraphiQL UI
}));

// Start the server
const PORT = 4002;
app.listen(PORT, () => {
    console.log(`GraphQL server running at http://localhost:${PORT}/query`);
});