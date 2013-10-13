<?php

if (!function_exists('index_collection')) {
  function get_list_value($item, $key) {
    if (is_callable($key)) {
      return $key($item);
    }
    return is_object($item) ? $item->{$key} : $item[$key];
  }

  function index_collection($collection, $value, $key = null) {
    if (is_object($collection)) {
      $collection = $collection->all();
    }
    return array_reduce($collection, function($result, $item) use ($value, $key) {
      $val = $item;
      if ($value) {
        $val = get_list_value($item, $value);
      }

      if ($key) {
        $result[get_list_value($item, $key)] = $val;
      } else {
        $result[] = $val;
      }

      return $result;
    }, array());
  }

  function index_collection_array($collection, $value, $key = null) {
    if (is_object($collection)) {
      $collection = $collection->all();
    }
    return array_reduce($collection, function($result, $item) use ($value, $key) {
      $val = $item;
      if ($value) {
        $val = get_list_value($item, $value);
      }

      if ($key) {
        $result[get_list_value($item, $key)][] = $val;
      } else {
        $result[] = $val;
      }

      return $result;
    }, array());
  }

  function logDB() {
    Event::listen("illuminate.query", function($query, $bindings, $time, $name) {
      echo "\n#" . $query . "\n";
      echo "___: " . json_encode($bindings) . "\n";
    });
  }
}

