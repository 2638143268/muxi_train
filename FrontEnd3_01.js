var checkIfInstanceOf = function (obj, classFunction) {
    if (obj == null || typeof classFunction !== 'function') {
        return false;
    }

    if (obj instanceof classFunction) {
        return true;
    }

    try {
        const wrapped = Object(obj);
        return wrapped instanceof classFunction;
    } catch (e) {
        return false;
    }
};

