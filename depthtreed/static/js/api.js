var Api = function(gateway) {
  this.gateway = gateway;
  this.rootNode = null;
  this.roots = [];
  this.nodes = {};
  this.clusters = [];
  this.nodeRelations = {};
};

Api.prototype.clear = function() {
  this.rootNode = null;
  this.roots = [];
  this.nodes = [];
  this.clusters = [];
  this.nodeRelations = {};
};

Api.prototype.nodeName = function(id) {
  return '' + id;
};

Api.prototype.newNode = function(v) {
  return {id: v.i, name: this.nodeName(v.i), max_depth: v.mxd || 0, min_depth: v.mnd || 0, children_count: v.n || 0};
};

Api.prototype.newCluster = function(v, clusterType) {
  var self = this;
  var cluster = {type: clusterType, id: 0, name: 'cluster', range: v.range, count: v.count};
  if (v.roots) {
    cluster.children = [];
    v.roots.forEach(function(r) {
      var node = self.newNode(r);
      self.nodes[node.name] = node;
      self.nodeRelations[node.name] = [];
      cluster.children.push(node);
    })
  }
  return cluster;
};

Api.prototype.newGenericNode = function(name) {
  return {id: 0, name: name, max_depth: 0, min_depth: 0, children_count: 10};
};

Api.prototype.findNodeById = function(id) {
  var nodeName = this.nodeName(id);
  return this.nodes[nodeName];
};

Api.prototype.getRootNode = function() {
  if (this.rootNode) {
    return this.rootNode;
  }
  return this.newGenericNode('generic');
};

Api.prototype.useDb = function(dbName) {
  this.db = dbName;
  this.clear();
};

Api.prototype.getDbs = function() {
  var self = this;
  return new Promise(function(resolve, reject) {
    fetch(self.gateway + '/db/list').then(function(resp) {
      if (resp.ok) {
        return resp.json();
      }
      throw new Error('Request failed!');
    }, function(error) {
      reject(error);
    }).then(function(json) {
      if (!json) {
        resolve(null);
        return;
      }
      if (json.code) {
        reject(new Error(json.message));
        return;
      }
      resolve(json);
    }).catch(function(error) {
      reject(error);
    });
  })
};

Api.prototype.getRoots = function() {
  var self = this;
  return new Promise(function(resolve, reject) {
    fetch(self.gateway + '/db/roots/' + self.db).then(function(resp) {
      if (resp.ok) {
        return resp.json();
      }
      throw new Error('Request failed!');
    }, function(error) {
      reject(error);
    }).then(function(json) {
      if (!json) {
        resolve(null);
        return;
      }
      if (json.code) {
        reject(new Error(json.message));
        return;
      }
      self.roots = [];
      json.forEach(function(v, idx, arr) {
        var node = self.newNode(v);
        self.roots.push(node);
        self.nodes[node.name] = node;
        self.nodeRelations[node.name] = [];
      });
      resolve(self.roots);
    }).catch(function(error) {
      reject(error);
    });
  })
};

Api.prototype.getTopChildren = function(depth, limit) {
  var self = this;
  return new Promise(function(resolve, reject) {
    fetch(self.gateway + '/db/top-children/' + self.db + '/' + depth + '/' + limit).then(function(resp) {
      if (resp.ok) {
        return resp.json();
      }
      throw new Error('Request failed!');
    }, function(error) {
      reject(error);
    }).then(function(json) {
      if (!json) {
        resolve(null);
        return;
      }
      if (json.code) {
        reject(new Error(json.message));
        return;
      }
      resolve(json);
    }).catch(function(error) {
      reject(error);
    });
  })
};

Api.prototype.getNode = function(nodeId) {
  var self = this;
  return new Promise(function(resolve, reject) {
    fetch(self.gateway + '/node/info/' + self.db + '/' + nodeId).then(function(resp) {
      if (resp.ok) {
        return resp.json();
      }
      throw new Error('Request failed!');
    }, function(error) {
      reject(error);
    }).then(function(json) {
      if (!json) {
        resolve(null);
        return;
      }
      if (json.code) {
        reject(new Error(json.message));
        return;
      }
      var node = self.newNode(json);
      self.nodes[node.name] = node;
      resolve(node);
    }).catch(function(error) {
      reject(error);
    });
  })
};

Api.prototype.getChildren = function(nodeId, depth) {
  var self = this;
  var nodeName = this.nodeName(nodeId);
  return new Promise(function(resolve, reject) {
    fetch(self.gateway + '/node/children/' + self.db + '/' + nodeId + '/' + depth).then(function(resp) {
      if (resp.ok) {
        return resp.json();
      }
      throw new Error('Request failed!');
    }, function(error) {
      reject(error);
    }).then(function(json) {
      if (!json) {
        resolve(null);
        return;
      }
      if (json.code) {
        reject(new Error(json.message));
        return;
      }
      var node = self.findNodeById(nodeId);
      node.children_count = json.count;
      self.flattenNodes(nodeId, json.nodes);
      resolve(true);
    }).catch(function(error) {
      reject(error);
    });
  })
};

Api.prototype.getDepthClusters = function(k) {
  var self = this;
  return new Promise(function(resolve, reject) {
    fetch(self.gateway + '/cluster/depth/' + self.db + '/' + k).then(function(resp) {
      if (resp.ok) {
        return resp.json();
      }
      throw new Error('Request failed!');
    }, function(error) {
      reject(error);
    }).then(function(json) {
      if (!json) {
        resolve(null);
        return;
      }
      if (json.code) {
        reject(new Error(json.message));
        return;
      }
      self.clusters = [];
      json.forEach(function(v, idx, arr) {
        var c = self.newCluster(v, 1);
        self.clusters.push(c);
      });
      resolve(self.clusters);
    }).catch(function(error) {
      reject(error);
    });
  })
};

Api.prototype.getChildrenClusters = function(depth, k) {
  var self = this;
  return new Promise(function(resolve, reject) {
    fetch(self.gateway + '/cluster/children/' + self.db + '/' + depth + '/' + k).then(function(resp) {
      if (resp.ok) {
        return resp.json();
      }
      throw new Error('Request failed!');
    }, function(error) {
      reject(error);
    }).then(function(json) {
      if (!json) {
        resolve(null);
        return;
      }
      if (json.code) {
        reject(new Error(json.message));
        return;
      }
      self.clusters = [];
      json.forEach(function(v, idx, arr) {
        var c = self.newCluster(v, 2);
        self.clusters.push(c);
      });
      resolve(self.clusters);
    }).catch(function(error) {
      reject(error);
    });
  })
};

Api.prototype.flattenNodes = function(rootId, nodes) {
  if (!nodes) {
    return;
  }
  var self = this;
  var nodeName = this.nodeName(rootId);
  nodes.forEach(function(v, idx, arr) {
    var node = self.newNode(v);
    self.nodes[node.name] = node;
    if (self.nodeRelations[nodeName]) {
      var found = self.nodeRelations[nodeName].find(function(id) {
        return id === node.id;
      });
      if (found) {
        return;
      }
    } else {
      self.nodeRelations[nodeName] = [];
    }
    self.nodeRelations[nodeName].push(node.id);
    if (v.c && v.c.length > 0) {
      self.flattenNodes(v.i, v.c);
    }
  });
};

Api.prototype.serializeClusters = function() {
  var root = this.newGenericNode('generic');
  if (!this.clusters || this.clusters.length === 0) {
    return root;
  }
  var self = this;
  root.children = [];
  this.clusters.forEach(function(cluster) {
    if (cluster.children) {
      cluster.children.forEach(function(node) {
        self.serialize(node);
      });
    }
    root.children.push(cluster);
  });
  return root;
};

Api.prototype.serialize = function(root, nodes) {
  if (nodes && nodes.length > 0) {
    var self = this;
    nodes.forEach(function(node) {
      node.children = self.serializeChildren(node);
      if (!root.children) {
        root.children = [];
      }
      root.children.push(node);
    })
    return root;
  }
  root.children = this.serializeChildren(root);
  return root;
};

Api.prototype.serializeChildren = function(root) {
  var nodeName = this.nodeName(root.id);
  if (!this.nodes[nodeName] || !this.nodeRelations[nodeName] || !this.nodeRelations[nodeName].length === 0) {
    return null;
  }
  var self = this;
  var children = [];
  this.nodeRelations[nodeName].forEach(function(id) {
    var childNode = self.findNodeById(id);
    if (childNode) {
      var nodes = self.serializeChildren(childNode);
      if (nodes && nodes.length > 0) {
        childNode.children = nodes;
      }
      children.push(childNode);
    }
  });
  return children;
};