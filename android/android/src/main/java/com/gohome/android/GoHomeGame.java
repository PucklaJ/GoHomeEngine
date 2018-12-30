package com.gohome.android;

import org.libsdl.app.SDLActivity;

public class GoHomeGame extends SDLActivity {

    @Override
    protected String[] getLibraries() {
        return new String[] {
            "gohome"
        };
    }

    @Override
    protected String getMainFunction() {
        return "SDL_main";
    }
}
