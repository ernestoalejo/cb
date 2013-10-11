'use strict';


// Use this DSL to delay some actions to the correct point
// in the chain of futures.
angular.scenario.dsl('defer', function() {
  return function(f) {
    return this.addFuture('deferred function', function(done) {
      f.call(this);
      done();
    });
  };
});
