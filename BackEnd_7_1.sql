# 数据库设计与SQL操作

## 1. 数据库设计
-- 书籍表
CREATE TABLE Books (
    ID VARCHAR(20) PRIMARY KEY,
    Title VARCHAR(100) NOT NULL,
    Author VARCHAR(100) NOT NULL
);

-- 库存表
CREATE TABLE Inventory (
    BookID VARCHAR(20) PRIMARY KEY,
    Stock INT NOT NULL,
    FOREIGN KEY (BookID) REFERENCES Books(ID)
);

-- 人员表
CREATE TABLE Persons (
    ID INT PRIMARY KEY,
    Name VARCHAR(50) NOT NULL
);

-- 喜好关系表（多对多关系）
CREATE TABLE ReadingPreferences (
    PersonID INT,
    BookID VARCHAR(20),
    PRIMARY KEY (PersonID, BookID),
    FOREIGN KEY (PersonID) REFERENCES Persons(ID),
    FOREIGN KEY (BookID) REFERENCES Books(ID)
);

## 2. 数据插入
-- 插入书籍数据
INSERT INTO Books (ID, Title, Author) VALUES
('go-away', 'the way to go', 'Ivo'),
('go-lang', 'Go语言圣经', 'Alan, Brian'),
('go-web', 'Go Web编程', 'Anonymous'),
('con-cur', 'Concurrency in Go', 'Katherine');

-- 插入库存数据
INSERT INTO Inventory (BookID, Stock) VALUES
('go-away', 20),
('go-lang', 17),
('go-web', 4),
('con-cur', 9);

-- 插入人员数据
INSERT INTO Persons (ID, Name) VALUES
(1, '小明'),
(2, '张三'),
(3, '翟曙');

-- 插入阅读喜好数据
INSERT INTO ReadingPreferences (PersonID, BookID) VALUES
(1, 'go-away'), (1, 'go-web'), (1, 'con-cur'),
(2, 'go-away'),
(3, 'go-web'), (3, 'con-cur');

## 3. 查询操作

### 1) 查询喜欢阅读#3的人
SELECT p.Name 
FROM Persons p
JOIN ReadingPreferences rp ON p.ID = rp.PersonID
WHERE rp.BookID = 'go-web';


### 2) 查询没有被喜欢阅读的书的信息(id,作者,标题,库存)
SELECT b.ID, b.Author, b.Title, i.Stock
FROM Books b
JOIN Inventory i ON b.ID = i.BookID
LEFT JOIN ReadingPreferences rp ON b.ID = rp.BookID
WHERE rp.BookID IS NULL;


### 3) 查询哪些人喜欢哪本书,列出人名和书名
SELECT p.Name, b.Title
FROM Persons p
JOIN ReadingPreferences rp ON p.ID = rp.PersonID
JOIN Books b ON rp.BookID = b.ID
ORDER BY p.Name, b.Title;
