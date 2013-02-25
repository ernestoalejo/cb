{{ if not .Exists }}'use strict';


var m = angular.module('{{ .Data.Module }}', []);
{{ end }}

m.factory('{{ .Data.Name }}', function() {
  return {};
});
