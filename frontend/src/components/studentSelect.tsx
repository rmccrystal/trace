import React, {Dispatch, SetStateAction, useEffect, useState} from "react";
import * as Api from "../api";
import {onCatch} from "./util";
import {ItemPredicate, ItemRenderer, ItemsEqualComparator, Select} from "@blueprintjs/select";
import {Button, MenuItem} from "@blueprintjs/core";

const StudentSelectElem = Select.ofType<Api.TraceStudent>();

export default function StudentSelect({onSelect}: { onSelect: (student: Api.TraceStudent | null) => void }) {
    let [students, setStudents] = useState<Api.TraceStudent[]>([]);

    let [activeStudent, setActiveStudent] = useState<Api.TraceStudent | null>(null);

    // useEffect(() => onSelect(activeLocation!), [onSelect, activeStudent]);

    useEffect(() => {
        Api.getStudents()
            .then(setStudents)
            .catch(onCatch)
    }, []);

    return <StudentSelectElem
        items={students}
        itemRenderer={renderStudent}
        itemPredicate={filterStudent}
        itemsEqual={areStudentsEqual}
        activeItem={activeStudent}
        onItemSelect={onSelect}
        noResults={<MenuItem disabled={true} text="No results."/>}
        className={"bp3-focus-disabled"}
    >
        <Button
            icon="person"
            rightIcon="caret-down"
            text={activeStudent?.name || "Select Student"}
        />
    </StudentSelectElem>
}

const renderStudent: ItemRenderer<Api.TraceStudent> = (student, {handleClick, modifiers, query}) => {
    return (
        <MenuItem
            key={student.id}
            onClick={handleClick}
            text={highlightText(student.name, query)}
            active={modifiers.active}
            disabled={modifiers.disabled}
        />
    );
}

const areStudentsEqual: ItemsEqualComparator<Api.TraceStudent> = ((itemA, itemB) => {
    return itemA.id === itemB.id
})

const filterStudent: ItemPredicate<Api.TraceStudent> = (query, student, _index, exactMatch) => {
    const normalizedTitle = student.name.toLowerCase();
    const normalizedQuery = query.toLowerCase();

    if (exactMatch) {
        return normalizedTitle === normalizedQuery;
    } else {
        return normalizedTitle.indexOf(normalizedQuery) >= 0;
    }
};

function escapeRegExpChars(text: string) {
    return text.replace(/([.*+?^=!:${}()|\[\]\/\\])/g, "\\$1");
}

function highlightText(text: string, query: string) {
    let lastIndex = 0;
    const words = query
        .split(/\s+/)
        .filter(word => word.length > 0)
        .map(escapeRegExpChars);
    if (words.length === 0) {
        return [text];
    }
    const regexp = new RegExp(words.join("|"), "gi");
    const tokens: React.ReactNode[] = [];
    while (true) {
        const match = regexp.exec(text);
        if (!match) {
            break;
        }
        const length = match[0].length;
        const before = text.slice(lastIndex, regexp.lastIndex - length);
        if (before.length > 0) {
            tokens.push(before);
        }
        lastIndex = regexp.lastIndex;
        tokens.push(<strong key={lastIndex}>{match[0]}</strong>);
    }
    const rest = text.slice(lastIndex);
    if (rest.length > 0) {
        tokens.push(rest);
    }
    return tokens;
}
