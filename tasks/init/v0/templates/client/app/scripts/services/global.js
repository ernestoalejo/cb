
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
    },

    cas: function(before, after) {
      if (msg_ == before)
        this.set(after);
    }
  };
});

m.factory('Selector', function() {
  var navbar_ = '';
  var dirtyNavbar_ = false;

  return {
    getNavbar: function() {
      return navbar_;
    },
    setNavbar: function(navbar) {
      navbar_ = navbar;
      dirtyNavbar_ = false;
    },

    setDirty: function() {
      dirtyNavbar_ = true;
    },
    clearDirty: function() {
      if (dirtyNavbar_)
        navbar_ = '';

      dirtyNavbar_ = false;
    },

    isDirty: function() {
      return dirtyNavbar_;
    },
    isNavbarDirty: function() {
      return dirtyNavbar_;
    }
  };
});
