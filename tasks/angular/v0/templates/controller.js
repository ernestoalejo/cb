{{ if not .Exists }}'use strict';


var m = angular.module('{{ .Data.Module }}', []);
{{ end }}

m.controller('{{ .Data.Name }}', function() {
  // empty
});
