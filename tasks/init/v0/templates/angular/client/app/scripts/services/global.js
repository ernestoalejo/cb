'use strict';


var m = angular.module('services.global', []);


m.factory('GlobalMsg', function($timeout) {
  var $globalMsg_ = $('#global-msg');
  var $notifications_ = $('#notifications');
  var last_;

  return {
    // Manage the loading message visibility
    showLoading: function() {
      $globalMsg_.show();
    },
    hideLoading: function() {
      $globalMsg_.hide();
    },
    isLoading: function() {
      return $globalMsg_.is(':visible');
    },

    // Create a new temporary message in the top-right corner of the page
    create: function(msg, type) {
      last_ = {msg: msg, type: type};

      type = type || 'info';
      if ($notifications_.notify) {  // for tests
        var not = $notifications_.notify({
          message: {text: msg},
          type: type,
          fadeOut: {
            enabled: true,
            delay: 10000
          }
        });
        not.$note.find('.close').click(function(e) {
          e.preventDefault();
        });
        not.show();
      }
    },

    // Useful for testing messages
    getLast: function() {
      return last_;
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
