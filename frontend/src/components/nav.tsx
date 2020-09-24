import React, {useEffect, useState} from "react";
import {Navbar, Button} from "@blueprintjs/core";
import "./nav.scss";
import LocationSelect from "./locationSelect";
import * as Api from "../api";
import { useGlobalState } from "../app";

export default function Nav() {
    let [location, setLocation] = useGlobalState('location')

    const onLocationSelect = (location: Api.Location) => {
        setLocation(location);
    }

    return <Navbar className="nav">
        <Navbar.Group align="center" className="items-center">
            <LocationSelect onSelect={onLocationSelect} />
        </Navbar.Group>
    </Navbar>
}