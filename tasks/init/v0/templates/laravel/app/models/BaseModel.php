<?php

class BaseModel extends Eloquent {

  protected $guarded = array();
  protected $cascadeDelete = array();

  public function delete() {
    Event::fire('model.delete', array($this));

    foreach ($this->cascadeDelete as $rel) {
      $items = $this->{$rel};
      if ($items instanceof \Illuminate\Database\Eloquent\Collection) {
        foreach ($items as $item) {
          $item->delete();
        }
      } else if ($items) {
        $items->delete();
      }
    }

    parent::delete();
  }

}
