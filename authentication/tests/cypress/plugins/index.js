// ***********************************************************
// This example plugins/index.js can be used to load plugins
//
// You can change the location of this file or turn off loading
// the plugins file with the 'pluginsFile' configuration option.
//
// You can read more here:
// https://on.cypress.io/plugins-guide
// ***********************************************************

// This function is called when a project is opened or re-opened (e.g. due to
// the project's config changing)
const MongoClient = require('mongodb').MongoClient;

module.exports = (on, config) => {
  // `on` is used to hook into various events Cypress emits
  // `config` is the resolved Cypress config
  on('task', {
    'db:teardown': () => {
      return new Promise((resolve) => {
        console.log("MONGO HOST: " + process.env.MONGO_HOST)
        MongoClient.connect(`mongodb://${process.env.MONGO_HOST || "mongo"}:27017`, (err, client) => {
          const db = client.db('prauxy')
          db.collection("users").deleteMany({})
          resolve("Done")
        })
      })
      // teardown
    },
    'db:seed': () => {
      return new Promise((resolve) => {
        MongoClient.connect(`mongodb://${process.env.MONGO_HOST || "mongo"}:27017`, (err, client) => {
          const db = client.db('prauxy')
          console.log(db)
          db.collection("users").insertOne({"username" : "test_user", 
                                            "email" : "demo@example.com", 
                                            "password" : "$2a$14$4X80jJ2XLaKj7h7unyvRAOPOY730GfJuvww3DTfvV9qQduG6i305." // abc123
                                          })

          db.collection("users").insertOne({"username" : "test_user_with_totp", 
                                            "email" : "demo-totp@example.com", 
                                            "password" : "$2a$14$4X80jJ2XLaKj7h7unyvRAOPOY730GfJuvww3DTfvV9qQduG6i305.", // abc123
                                            "multifactor" : [ 
                                              {
                                               "type" : "totp", 
                                               "secret" : "MZZDVPZZTWG4WZYOBUTJXMIOMP37AKZ2" 
                                              } 
                                            ]
                                          })

          resolve("Done")
        })
      })
    }
  })
}
