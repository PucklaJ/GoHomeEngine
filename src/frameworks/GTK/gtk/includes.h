#ifndef GTK_GO_INCLUDES
#define GTK_GO_INCLUDES

#include <gtk/gtk.h>
#include <stdlib.h>

typedef struct {
    int id;
    char* signal;
} ButtonSignalUserData;

extern GtkWindow* Window;
extern GtkGLArea* GLarea;

extern char* ErrorString;

extern void initialise(int args,char** argv);

extern int createWindow(unsigned int width, unsigned int height, const char* title);
extern GtkWindow* createWindowObject();
extern void configureWindowParameters(GtkWindow* window,unsigned int width, unsigned int height, const char* title);
extern void connectWindowSignals(GtkWindow* window);
extern void createGLArea();
extern void configureGLArea(GtkGLArea* area);
extern void addGLAreaToWindow();
extern void addGLAreaToContainer(GtkContainer* container);
extern void windowSetSize(float width, float height);
extern void windowGetSize(float* width, float* height);
extern void windowHideCursor();
extern void windowShowCursor();
extern void windowDisableCursor();
extern int windowCursorShown();
extern int windowCursorHidden();
extern int windowCursorDisabled();

extern void loop();


extern GtkContainer* boxToContainer(GtkBox* box);
extern GtkContainer* buttonToContainer(GtkButton* button);
extern GtkWidget* boxToWidget(GtkBox* box);
extern GtkWidget* glareaToWidget(GtkGLArea* area);
extern GtkWidget* buttonToWidget(GtkButton* button);
extern GtkContainer* windowToContainer(GtkWindow* window);
extern GtkBox* widgetToBox(GtkWidget* widget);
extern GtkButton* widgetToButton(GtkWidget* widget);
extern GtkWidget* gobjectToWidget(GObject* object);
extern GObject* widgetToGObject(GtkWidget* widget);
extern GtkWindow* widgetToWindow(GtkWidget* widget);
extern GtkWidget* gpointerToWidget(gpointer data);
extern GtkContainer* widgetToContainer(GtkWidget* widget);
extern GtkGrid* widgetToGrid(GtkWidget* widget);
extern GtkWidget* windowToWidget(GtkWindow* window);
extern GtkGLArea* gobjectToGLArea(GObject* object);

extern void signalConnectButton(GtkButton* button,char* signal, int id);

#endif