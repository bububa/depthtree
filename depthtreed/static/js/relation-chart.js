var RelationChart = function(selector, w, h) {
  this.w = w;
  this.h = h;
  this.updateNodeFunc = null;

  d3.select(selector).selectAll("svg").remove();

  this.svg = d3.select(selector).append("svg:svg")
    .attr('width', w)
    .attr('height', h);

  this.svg.append("svg:rect")
    .style("stroke", "#999")
    .style("fill", "#000")
    .attr('width', w)
    .attr('height', h);

  this.force = d3.layout.force()
    .on("tick", this.tick.bind(this))
    .charge(function(d) { return d._children ? -d.size / 100 : -40; })
    .linkDistance(function(d) { return d.target._children ? 80 : 25; })
    .size([h, w]);
};

RelationChart.prototype.update = function(json) {
  if (json) this.json = json;

  this.json.fixed = true;
  this.json.x = this.w / 2;
  this.json.y = this.h / 2;

  var nodes = this.flatten(this.json);
  var links = d3.layout.tree().links(nodes);
  var total = nodes.length || 1;

  // remove existing text (will readd it afterwards to be sure it's on top)
  this.svg.selectAll("text").remove();

  // Restart the force layout
  this.force
    .gravity(Math.atan(total / 50) / Math.PI * 0.4)
    .nodes(nodes)
    .links(links)
    .start();

  // Update the links
  this.link = this.svg.selectAll("line.link")
    .data(links, function(d) { return d.target.name; });

  // Enter any new links
  this.link.enter().insert("svg:line", ".node")
    .attr("class", "link")
    .attr("x1", function(d) { return d.source.x; })
    .attr("y1", function(d) { return d.source.y; })
    .attr("x2", function(d) { return d.target.x; })
    .attr("y2", function(d) { return d.target.y; });

  // Exit any old links.
  this.link.exit().remove();

  // Update the nodes
  this.node = this.svg.selectAll("circle.node")
    .data(nodes, function(d) { return d.name; })
    .classed("collapsed", function(d) { return d._children ? 1 : 0; });

  this.node.transition()
    // .attr("r", function(d) { return d.children ? 3.5 : Math.pow(d.size, 2/5) || 1; });
    .attr("r", function(d) { return d.id === 0 ? 3.5 : Math.pow(d.size, 2/5) || 1; });

  // Enter any new nodes
  this.node.enter().append('svg:circle')
    .attr("class", "node")
    .classed('directory', function(d) { return (d._children || d.children) ? 1 : 0; })
    // .attr("r", function(d) { return d.children ? 3.5 : Math.pow(d.size, 2/5) || 1; })
    .attr("r", function(d) { return d.id === 0 ? 3.5 : Math.pow(d.size, 2/5) || 1; })
    .style("fill", function color(d) {
      if (d.id === 0) {
        return "rgba(255, 255, 255, 1)";
      }
      return "hsl(" + parseInt(360 / total * d.id, 10) + ",90%,70%)";
    })
    .call(this.force.drag)
    .on("click", this.click.bind(this))
    .on("mouseover", this.mouseover.bind(this))
    .on("mouseout", this.mouseout.bind(this));

  // Exit any old nodes
  this.node.exit().remove();

  this.text = this.svg.append('svg:text')
    .attr('class', 'nodetext')
    .attr('dy', 0)
    .attr('dx', 0)
    .attr('text-anchor', 'middle');

  return this;
};

RelationChart.prototype.flatten = function(root) {
  var nodes = [];

  function recurse(node) {
    if (node.children) {
      node.children.forEach(function(n) {
        n.size = recurse(n);
      })
      /*
      node.size = node.children.reduce(function(p, v) {
        return p + recurse(v);
      }, 0);
      */
    }
    // if (node.id) node.id = ++i;
    nodes.push(node);
    node.size = node.children_count || 1;
    return node.size;
  }

  root.size = recurse(root);
  return nodes;
};

RelationChart.prototype.click = function(d) {
  // Toggle children on click.
  if (d.children && d.children.length > 0) {
    d._children = d.children;
    d.children = null;
    this.update();
  } else if (d._children && d._children.length > 0) {
    d.children = d._children;
    d._children = null;
    this.update();
  } else if (d.id > 0 && this.updateNodeFunc) {
    this.updateNodeFunc(d);
  }
};

RelationChart.prototype.mouseover = function(d) {
  var txt = '[' + d.name + '] maxDepth:' + d.max_depth + ' children:' + d.children_count;
  if (d.type && d.type === 1) {
    txt = '[depth cluster] range:' + d.range[0] + '-' + d.range[1] + ' count:' + d.count;
  } else if (d.type && d.type === 2) {
    txt = '[children cluster] range:' + d.range[0] + '-' + d.range[1] + ' count:' + d.count;
  }
  this.text.attr('transform', 'translate(' + d.x + ',' + (d.y - 5 - (d.children ? 3.5 : Math.sqrt(d.size) / 2)) + ')')
    .text(txt)
    .style('display', null);
};

RelationChart.prototype.mouseout = function(d) {
  this.text.style('display', 'none');
};

RelationChart.prototype.tick = function() {
  var h = this.h;
  var w = this.w;
  this.link.attr("x1", function(d) { return d.source.x; })
    .attr("y1", function(d) { return d.source.y; })
    .attr("x2", function(d) { return d.target.x; })
    .attr("y2", function(d) { return d.target.y; });

  this.node.attr("transform", function(d) {
    return "translate(" + Math.max(5, Math.min(w - 5, d.x)) + "," + Math.max(5, Math.min(h - 5, d.y)) + ")";
  });
};

RelationChart.prototype.cleanup = function() {
  this.update([]);
  this.force.stop();
};