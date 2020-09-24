import React, {useEffect, useState} from "react";
import {Navbar, Button} from "@blueprintjs/core";
import "./nav.scss";
import LocationSelect from "./locationSelect";
import * as Api from "../api";
import { useGlobalState } from "../app";

export default function Nav() {
    let [location, setLocation] = useGlobalState('location')

    let [dark, setDark] = useGlobalState('dark');
    const onToggleDark = () => {
        setDark(!dark);
    }

    const onLocationSelect = (location: Api.Location) => {
        setLocation(location);
    }

    return <Navbar fixedToTop>
        <Navbar.Group align="left">
            <LocationSelect onSelect={onLocationSelect} />
        </Navbar.Group>
        <Navbar.Group align="right">
            <Button icon="contrast" minimal onClick={onToggleDark}/>
        </Navbar.Group>
    </Navbar>
}