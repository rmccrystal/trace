import React, {Dispatch, SetStateAction, useEffect, useState} from "react";
import * as Api from "../api";
import {onCatch} from "./util";
import {ItemPredicate, ItemRenderer, ItemsEqualComparator, Select} from "@blueprintjs/select";
import {Button, MenuItem, Navbar} from "@blueprintjs/core";

const LocationSelectElem = Select.ofType<Api.Location>();

export default function LocationSelect({onSelect}: { onSelect: (location: Api.Location) => void }) {
    let [locations, setLocations] = useState<Api.Location[]>([]);
    let [activeLocation, setActiveLocation] = useState<Api.Location | null>(null)

    useEffect(() => onSelect(activeLocation!), [onSelect, activeLocation]);

    useEffect(() => {
        Api.getLocations()
            .then((locs) => {
                setLocations(locs);
                setActiveLocation(locs[0]);
            })
            .catch(onCatch)
    }, []);

    return <LocationSelectElem
        items={locations}
        itemRenderer={renderLocation}
        itemPredicate={filterLocation}
        itemsEqual={areLocationsEqual}
        activeItem={activeLocation}
        createNewItemRenderer={renderCreateLocationOption}
        onItemSelect={setActiveLocation}
        noResults={<MenuItem disabled={true} text="No results." />}
        className={"bp3-focus-disabled"}
    >
        <Button
            icon="office"
            rightIcon="caret-down"
            text={activeLocation?.name}
        />
    </LocationSelectElem>
}

// We have to use a generator for this so we can update the location list
function createLocationGenerator(setLocations: Dispatch<SetStateAction<Api.Location[]>>): (name: string) => Api.Location {
    return name => {
        const newLocation = {name: "test", id: "hello", timeout: 1};
        setLocations(prevState => [...prevState, newLocation]);
        return newLocation
    }
}

const renderLocation: ItemRenderer<Api.Location> = (location, { handleClick, modifiers, query }) => {
    return (
        <MenuItem
            key={location.id}
            onClick={handleClick}
            text={highlightText(location.name, query)}
            active={modifiers.active}
            disabled={modifiers.disabled}
        />
    );
}

export const renderCreateLocationOption = (
    query: string,
    active: boolean,
    handleClick: React.MouseEventHandler<HTMLElement>,
) => (
    <MenuItem
        icon="add"
        text={`Create "${query}"`}
        active={active}
        onClick={handleClick}
        shouldDismissPopover={false}
    />
);

const areLocationsEqual: ItemsEqualComparator<Api.Location> = ((itemA, itemB) => {
    return itemA.id === itemB.id
})

const filterLocation: ItemPredicate<Api.Location> = (query, location, _index, exactMatch) => {
    const normalizedTitle = location.name.toLowerCase();
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
