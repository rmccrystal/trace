import React from "react";
import {Button, IconName, Navbar} from "@blueprintjs/core";
import LocationSelect from "./locationSelect";
import * as Api from "../api";
import {Link} from "react-router-dom";

export interface NavProps {
    location: Api.TraceLocation,
    setLocation: (location: Api.TraceLocation) => void,
    onToggleDark: () => void,
}

export default function Nav({location, setLocation, onToggleDark}: NavProps) {
    const NavLink = ({icon, title, to}: { icon: IconName, title: string, to: string }) => <Link to={to}>
        <Button minimal className={"mx-1"} icon={icon} text={title} />
    </Link>

    return <Navbar>
        <Navbar.Group align="left">
            <LocationSelect activeLocation={location} onSelect={setLocation}/>
            <Navbar.Divider/>

            <NavLink icon="align-justify" title="Scan" to="/" />
            <NavLink icon="dashboard" title={`${location.name} Dashboard`} to="/dashboard" />
            <NavLink icon="people" title="Manage Students" to="/students" />
            <NavLink icon="graph" title="Contact Tracing" to="/trace" />
        </Navbar.Group>
        <Navbar.Group align="right">
            <Button icon="contrast" minimal onClick={onToggleDark}/>
        </Navbar.Group>
    </Navbar>
}