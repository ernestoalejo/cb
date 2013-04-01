'use strict';


var m = angular.module('services.global', []);


m.factory('GlobalMsg', function($timeout) {
  var msg_ = '', tm = null, type_ = 'success';
  
  return {
    set: function(msg, type) {
      msg_ = msg;
      type_ = type ? type : 'success';
      this.cleanTimer();
    },

    cleanTimer: function(msg, timer) {
      if (tm) {
        $timeout.cancel(tm);
        tm = null;
      }
    },

    setTemp: function(msg, type) {
      this.cleanTimer();
      this.set(msg, type);

      var that = this;
      tm = $timeout(function() {
        that.set('');
      }, 10000);
    },

    get: function() {
      return msg_;
    },

    getClass: function() {
      return 'label-' + type_;
    }
  };
});


m.factory('ErrorRegister', function() {
  var error_ = null;

  return {
    clean: function() {
      error_ = null;
    },

    set: function(error) {
      error_ = error;
    },

    isNull: function() {
      return error_ === null;
    }
  };
});
