import React, {Dispatch, SetStateAction, useCallback, useEffect, useState} from 'react';
import './app.scss';
import Scan from './components/scan';
import Nav from "./components/nav";
import {Card, FocusStyleManager, Spinner} from "@blueprintjs/core";
import {Switch, Route, Redirect, BrowserRouter} from "react-router-dom";
import {getLocations, TraceLocation} from "./api";
import useLocalStorage, {onCatchPrefix} from "./components/util";
import StudentList from "./components/studentList";
import CurrentlyInLocation from "./components/currentlyInLocation";
import NewLocationPrompt from "./components/newLocationPrompt";

function App() {
    useEffect(() => {
        FocusStyleManager.onlyShowFocusOnTabs();
    }, [])

    let [dark, setDark] = useLocalStorage('dark', false);
    let [location, setLocation] = useState<TraceLocation | null>(null);
    let [locationID, setLocationID] = useLocalStorage('location', '');
    let [newLocationPrompt, setNewLocationPrompt] = useState(false);

    useEffect(() => {
        if (location) {
            setLocationID(location.id)
        }
    }, [location, setLocationID])

    const updateLocations = useCallback(() => {
        getLocations()
            .then(locations => {
                if (locations.length === 0) {
                    setNewLocationPrompt(true);
                } else {
                    setLocation(locations.find(value => value.id === locationID) || locations[0])
                }
            })
            .catch(onCatchPrefix("Error getting list of locations: "));
    }, [locationID])

    useEffect(() => {
        updateLocations();
    }, [updateLocations])

    let innerElement;

    if (newLocationPrompt && location === null) {
        innerElement = <NewLocationPrompt submitCallback={setLocation} elevation={1} className="p-10 m-auto"/>
    } else if (location === null) {
        innerElement = <Spinner className="m-auto absolute inset-0"/>
    } else {
        innerElement = <><Nav setLocation={setLocation} location={location} onToggleDark={() => setDark(!dark)}/>
            <Switch>
                <Route exact path="/scan" component={() => <Scan location={location!} elevation={1}/>}/>
                <Route path="/dashboard"
                       component={() => <CurrentlyInLocation location={location!} elevation={1}/>}/>
                <Route path="/students" component={() => <StudentList elevation={1}/>}/>
                <Redirect from="*" to="/scan"/>
            </Switch>
            <div className="mt-auto text-center mb-3 bp3-text-muted">
                Created by Ryan McCrystal | <a href="https://github.com/rmccrystal/trace" rel="noopener noreferrer"
                                               target="_blank">github</a>
            </div>
        </>
    }

    return (
        <BrowserRouter>
            <div className={`app flex flex-col items-center content-center ${dark ? "bp3-dark" : "bp3-light"}`}>
                {innerElement}
            </div>
        </BrowserRouter>
    );
}

export default App;