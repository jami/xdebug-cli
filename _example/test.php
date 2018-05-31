<?php

function Foo() {

    for ($i=0; $i < 100; $i++) {
        echo "Baaaaaa " . $i . PHP_EOL;
    }

    return "Bar";
}

class Bar {
    public function toString() {
        return Foo();
    }
}

$a = [
    "foo1" => "bar",
    "foo2" => 1,
    "foo3" => true
];

$bar = new Bar();
echo $bar->toString(); 