import React, {useCallback, useEffect, useState} from 'react';
import './app.scss';
import Scan from './components/scan';
import Nav from "./components/nav";
import {FocusStyleManager, Spinner} from "@blueprintjs/core";
import {BrowserRouter, Redirect, Route, Switch} from "react-router-dom";
import {getLocations, TraceLocation} from "./api";
import useLocalStorage, {onCatchPrefix} from "./components/util";
import ManageStudents from "./components/manageStudents";
import CurrentlyInLocation from "./components/currentlyInLocation";
import NewLocationPrompt from "./components/newLocationPrompt";
import VisitedLocationToday from "./components/visitedLocationToday";

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

    useEffect(() => {
        if (dark) {
            document.body.classList.remove("bp3-light");
            document.body.classList.add("bp3-dark");
        } else {
            document.body.classList.remove("bp3-dark");
            document.body.classList.add("bp3-light");
        }
    }, [dark])

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
                <Route path="/dashboard">
                    <CurrentlyInLocation location={location!} elevation={1}/>
                    <VisitedLocationToday location={location!} elevation={1}/>
                </Route>
                <Route path="/students" component={() => <ManageStudents elevation={1}/>}/>
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
            <div className={`app flex flex-col items-center content-center`}>
                {innerElement}
            </div>
        </BrowserRouter>
    );
}

export default App;