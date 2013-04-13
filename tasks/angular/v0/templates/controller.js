{{ if not .Exists }}'use strict';


var m = angular.module('controllers.{{ .Data.Module }}', []);
{{ end }}

m.controller('{{ .Data.Name }}', function() {
  // empty
});
