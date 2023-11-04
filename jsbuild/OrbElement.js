import { Orb } from "@memgraph/orb";

class OrbGraph extends HTMLElement {
  constructor() {
    super()
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

    const orb = new Orb(this);
    orb.data.setup({ nodes, edges });
    orb.view.render(() => {
      orb.view.recenter();
    });
  }

}
customElements.define('orb-graph', OrbGraph)
