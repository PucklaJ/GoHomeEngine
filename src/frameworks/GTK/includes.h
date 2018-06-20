#ifndef GTK_GO_INCLUDES
#define GTK_GO_INCLUDES

#include <gtk/gtk.h>
#include <stdlib.h>

extern GtkWindow* Window;
extern GtkGLArea* GLarea;

extern char* ErrorString;

extern void initialise(int args,char** argv);

extern int createWindow(unsigned int width, unsigned int height, const char* title);
extern void windowSetSize(float width, float height);
extern void windowGetSize(float* width, float* height);
extern void windowHideCursor();
extern void windowShowCursor();
extern void windowDisableCursor();
extern int windowCursorShown();
extern int windowCursorHidden();
extern int windowCursorDisabled();

extern void loop();

#endif