console.log('1');

setTimeout(() => {
    console.log('2');
    new Promise(resolve => {
        console.log('3');
        resolve();
    }).then(() => {
        console.log('4');
    });
}, 0);

new Promise((resolve) => {
    console.log('5');
    resolve();
}).then(() => {
    console.log('6');
    return new Promise(resolve => {
        console.log('7');
        resolve();
    }).then(() => {
        console.log('8');
    });
}).then(() => {
    console.log('9');
});

Promise.resolve().then(() => {
    console.log('10');
    setTimeout(() => {
        console.log('11');
    }, 0);
    Promise.resolve().then(() => {
        console.log('12');
    });
});

console.log('13');

setTimeout(() => {
    console.log('14');
}, 0);

console.log('15');

// 推测执行顺序：
// 同步任务阶段（直接执行）：
// console.log('1') → 输出 1
// console.log('5') → 输出 5
// console.log('13') → 输出 13
// console.log('15') → 输出 15

// 微任务阶段（Promise.then）：
// console.log('6')（来自 new Promise(...).then(...)）
// console.log('7')（在 then 中创建的 Promise）
// console.log('8')（上面 Promise 的 then）
// console.log('9')（上面链式 then）
// console.log('10')
// console.log('12')

// 宏任务阶段（setTimeout）：
// 先执行第一个 setTimeout：
// console.log('2')
// console.log('3')
// console.log('4')（Promise.then）
// 然后执行 console.log('14')
// 最后是 console.log('11')（出现在 10 中的 setTimeout）

