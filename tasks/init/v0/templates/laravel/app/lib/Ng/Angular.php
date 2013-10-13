<?php namespace Ng;


class Angular {

  // To emmit Angular bindings inside Blade templates
  public static function bind($statement) {
    return  '{{' . $statement . '}}';
  }

}
