import { Orb } from "@memgraph/orb";

class OrbGraph extends HTMLElement {
  ws = null
  orb = null
  nodes = []
  edges = []
  constructor() {
    super()
    const self = this
    this.ws = new WebSocket("ws://localhost:1323/ws")
    this.nodes = []
    this.edges = []
    this.ws.onopen = function () {
      console.log("Client Opened")
    }
    this.ws.onclose = function () {
      console.log("Client Closed")
    }
    this.ws.onmessage = function (event) {
      console.log(event)
      const data = JSON.parse(event.data)
      switch (data.action) {
        case "CREATE":
          self.nodes.push({
            id: data.result.in,
            label: data.result.in
          })
          self.nodes.push({
            id: data.result.out,
            label: data.result.out
          })
          self.edges.push({
            id: data.result.id,
            start: data.result.in,
            end: data.result.out
          })
          break;
        case "UPDATE":

          break;
        case "DELETE":
          break;
        default:
          break;
      }
      self.updateGraph()
    }
  }

  connectedCallback() {
    const nodes = [
      { id: 1, label: "Node 1" },
      { id: 2, label: "Node 2" },
      { id: 3, label: "Node 3" },
    ];

    const edges = [
      { id: 1, start: 1, end: 2 },
      { id: 2, start: 2, end: 3 },
      { id: 3, start: 3, end: 1 },
    ];

    this.orb = new Orb(this);
    this.orb.data.setup({ nodes, edges });
    this.orb.view.render(() => {
      this.orb.view.recenter();
    });
  }

  updateGraph() {
    this.orb.data.setup({ nodes: this.nodes, edges: this.edges });
    this.orb.view.render(() => {
      this.orb.view.recenter();
    });
  }
}
customElements.define('orb-graph', OrbGraph)
