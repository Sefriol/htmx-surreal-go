const path = require("path");

module.exports = {
  entry: "./OrbElement.js",
  output: {
    filename: "bundle.js",
    path: path.resolve(__dirname, "../dist"),
  },
  module: {
    rules: [
      {
        test: /\.worker\.js$/,
        use: {
          loader: 'worker-loader',
          options: {
            filename: 'orb.worker.js',
          },
        }
      }
    ]
  },
  resolve: {
    alias: {
      "@memgraph/orb": path.resolve(__dirname, "node_modules/@memgraph/orb/dist/browser/orb.js"),
    },
  },
};
