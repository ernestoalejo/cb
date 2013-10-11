'use strict';


var m = angular.module('services.global', []);


m.factory('GlobalMsg', function($timeout) {
  var $notifications = $('#notifications');
  var loading = false, totalLoading = 0, cancelTimeout;

  return {
    // Manage the loading message visibility
    showLoading: function() {
      totalLoading++;
      this.timeout_();
    },
    hideLoading: function() {
      totalLoading--;
      totalLoading = Math.max(totalLoading, 0);
      this.timeout_();
    },
    forceHideLoading: function() {
      loading = false;
      totalLoading = 0;
    },
    isLoading: function() {
      return loading;
    },
    timeout_: function() {
      if (cancelTimeout) {
        $timeout.cancel(cancelTimeout);
      }
      cancelTimeout = $timeout(function() {
        loading = (totalLoading > 0);
        cancelTimeout = null;
      }, 1000);
    },

    // Create a new message in the top-right corner of the page
    create: function(msg, type, permanent) {
      type = type || 'info';
      if ($notifications.notify) {  // for tests
        var not = $notifications.notify({
          message: {text: msg},
          type: type,
          fadeOut: {
            enabled: !permanent,
            delay: 10000
          }
        });
        not.$note.find('.close').click(function(e) {
          e.preventDefault();
        });
        not.show();
      }
    }
  };
});


m.factory('ErrorRegister', function() {
  var error = null;

  return {
    clean: function() {
      error = null;
    },

    set: function(err) {
      error = err;
    },

    isNull: function() {
      return error === null;
    }
  };
});
