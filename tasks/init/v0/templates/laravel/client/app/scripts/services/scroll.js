'use strict';


var m = angular.module('services.scroll', []);


m.factory('ScrollTo', function($timeout) {
  var $scrollTarget = $('html,body');

  function scroll($elm) {
    if (typeof $elm == 'string') {
      $elm = $('#' + $elm);
    }

    if ($elm && $elm.size() > 0) {
      var offset = $elm.offset();
      offset.left -= 102;
      offset.top -= 82;

      $scrollTarget.animate({
        scrollTop: offset.top,
        scrollLeft: offset.left
      });
    } else {
      $scrollTarget.animate({
        scrollTop: 0,
        scrollLeft: 0
      });
    }
  }

  scroll.formError = function() {
    $timeout(function() {
      var err = $('.form-group.has-error:first');
      if (err.length > 0) {
        scroll(err);
      }
    }, 0);
  };

  return scroll;
});

