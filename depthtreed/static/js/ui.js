var NodeViewMode = 'node';
var DepthClusterViewMode = 'depthCluster';
var ChildrenClusterViewMode= 'childrenCluster';

var DEFAULT_DEPTH = 3;
var MAX_DEPTH = 10;
var DEFAULT_K = 10;

var Tasks = function(amount, callback) {
  this.amount = amount;
  this.callback = callback;
  this.results = [];
}

Tasks.prototype.done = function(res) {
  this.amount --;
  this.results.push(res);
  if (this.amount <= 0 && this.callback) {
    this.callback(this.results);
  }
}

var UI = function(chart, api) {
  this.chart = chart;
  this.api = api;
  this.depth = DEFAULT_DEPTH;
  this.k = DEFAULT_K;
  this.viewMode = NodeViewMode;
  var self = this;
  this.chart.updateNodeFunc = function(node) {
    self.updateNodeView(node.id)
  };
}

UI.prototype.updateNodeView = function(nodeId) {
  var self = this;
  var depth = this.depth;
  var rootNode = this.api.getRootNode();
  var roots = this.api.roots;
  if (rootNode.id === 0 && !nodeId) {
    this.api.nodeRelations = {};
    if (this.viewMode === NodeViewMode) {
      this.api.getRoots().then(function(rootNodes) {
        if (!rootNodes) {
          self.chart.update(null);
          return;
        }
        var tasks = new Tasks(rootNodes.length, function() {
          var json = self.api.serialize(rootNode, rootNodes);
          self.chart.update(json);
        });

        rootNodes.forEach(function(root) {
          self.api.getChildren(root.id, depth).then(function() {
            tasks.done(null);
          }, function(err) {
            tasks.done(err);
            alert(err);
          });
        });
      }, function(err) {
        alert(err);
      });
    } else if (this.viewMode === DepthClusterViewMode) {
      this.api.getDepthClusters(this.k).then(function(clusters) {
        if (!clusters) {
          self.chart.update(null);
          return;
        }
        var json = self.api.serializeClusters();
        self.chart.update(json);
      }, function(err) {
        alert(err);
      });
    } else if (this.viewMode === ChildrenClusterViewMode) {
      this.api.getChildrenClusters(this.depth, this.k).then(function(clusters) {
        if (!clusters) {
          self.chart.update(null);
          return;
        }
        var json = self.api.serializeClusters();
        self.chart.update(json);
      }, function(err) {
        alert(err);
      });
    }
  } else {
    nodeId = !nodeId && rootNode.id  > 0 ? rootNode.id : nodeId;
    self.api.getChildren(nodeId, depth).then(function() {
      var json = null;
      if (self.viewMode === NodeViewMode) {
        json = self.api.serialize(rootNode, rootNode.id === 0 ? roots : null);
      } else if (self.viewMode === DepthClusterViewMode || self.viewMode === ChildrenClusterViewMode) {
        json = self.api.serializeClusters();
      }
      self.chart.update(json);
    }, function(err) {
      alert(err);
    });
  }
};

UI.prototype.updateTopChildrenView = function(ol, limit) {
  this.api.getTopChildren(this.depth, limit).then(function(nodes) {
    var wrapper = document.querySelector(ol);
    if (!wrapper) {
      return;
    }
    wrapper.childNodes.forEach(function(e) { e.remove(); });
    nodes.forEach(function(node) {
      var li = document.createElement('li');
      li.innerText = '[' + node.i + '] children:' + node.n;
      li.dataset.id = node.i;
      wrapper.appendChild(li);
    });
  }, function(err) {
    alert(err);
  });
};

UI.prototype.changeDb = function(dbname) {
  this.api.useDb(dbname);
  this.updateNodeView();
};

UI.prototype.setRootNode = function(nodeId) {
  var self = this;
  this.viewMode = NodeViewMode;
  this.api.nodeRelations = {};
  if (nodeId === 0) {
    this.api.rootNode = null;
    this.updateNodeView();
    return;
  }
  this.api.getNode(nodeId).then(function(node) {
    self.api.rootNode = node;
    self.updateNodeView();
  })
};

UI.prototype.updateViewMode = function(viewMode) {
  this.viewMode = viewMode;
  this.api.clear();
  this.updateNodeView();
  return;
};

UI.prototype.updateDepth = function(depth) {
  this.depth = depth === 0 || depth > MAX_DEPTH ? DEFAULT_DEPTH : depth;
  this.updateNodeView(null);
}