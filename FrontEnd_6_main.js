import { double } from './FrontEnd_6_utils.js';

// 1. 箭头函数 multiply
const multiply = (a, b) => a * b;
console.log("multiply(3, 4):", multiply(3, 4)); // 输出：12

// 2. 解构数组
const [first, , third] = [1, 2, 3];
console.log("first:", first);   // 输出：1
console.log("third:", third);   // 输出：3

// 3. 使用导入的 double 函数
console.log("double(6):", double(6)); // 输出：12

// 4. 类与继承
class Animal {
    constructor(name) {
        this.name = name;
    }

    speak() {
        console.log(`${this.name} makes a sound.`);
    }
}

class Dog extends Animal {
    bark() {
        console.log("Woof!");
    }
}

const myDog = new Dog("Buddy");
myDog.speak(); // 输出：Buddy makes a sound.
myDog.bark();  // 输出：Woof!
