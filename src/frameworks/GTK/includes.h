#ifndef GTK_GO_INCLUDES
#define GTK_GO_INCLUDES

#include <gtk/gtk.h>

extern GtkWindow* Window;
extern GtkGLArea* GLarea;

extern char* ErrorString;

extern void initialise(int args,char** argv);

extern int createWindow(unsigned int width, unsigned int height, const char* title);
extern void windowGetSize(float* width, float* height);

extern void loop();

#endif