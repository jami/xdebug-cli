<?php

function Foo() {
    return "Bar";
}

class Bar {
    public function toString() {
        return Foo();
    }
}


$bar = new Bar();
echo $bar->toString(); 