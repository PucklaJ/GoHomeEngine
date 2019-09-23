#ifndef GTK_GO_INCLUDES
#define GTK_GO_INCLUDES

#include <gtk/gtk.h>
#include <stdlib.h>

typedef struct
{
    int id;
    char *signal;
} ButtonSignalUserData;

typedef struct
{
    char *name;
    char *signal;
} WidgetSignalUserData;

extern GtkWindow *Window;
extern GtkGLArea *GLarea;
extern char *ErrorString;

extern void initialise(int args, char **argv);

extern int createWindow(int width, int height, const char *title);
extern GtkWindow *createWindowObject();
extern void configureWindowParameters(GtkWindow *window, int width, int height, const char *title);
extern void connectWindowSignals(GtkWindow *window);
extern void createGLArea();
extern void configureGLArea(GtkGLArea *area);
extern void addGLAreaToWindow();
extern void addGLAreaToContainer(GtkContainer *container);
extern void widgetGetSize(GtkWidget *widget, gint *width, gint *height);
extern void windowSetSize(float width, float height);
extern void windowGetSize(float *width, float *height);
extern void windowHideCursor();
extern void windowShowCursor();
extern void windowDisableCursor();
extern int windowCursorShown();
extern int windowCursorHidden();
extern int windowCursorDisabled();

extern void loop();

extern GtkContainer *boxToContainer(GtkBox *box);
extern GtkContainer *buttonToContainer(GtkButton *button);
extern GtkWidget *boxToWidget(GtkBox *box);
extern GtkWidget *glareaToWidget(GtkGLArea *area);
extern GtkWidget *buttonToWidget(GtkButton *button);
extern GtkContainer *windowToContainer(GtkWindow *window);
extern GtkBox *widgetToBox(GtkWidget *widget);
extern GtkButton *widgetToButton(GtkWidget *widget);
extern GtkWidget *gobjectToWidget(GObject *object);
extern GObject *widgetToGObject(GtkWidget *widget);
extern GtkWindow *widgetToWindow(GtkWidget *widget);
extern GtkWidget *gpointerToWidget(gpointer data);
extern GtkContainer *widgetToContainer(GtkWidget *widget);
extern GtkGrid *widgetToGrid(GtkWidget *widget);
extern GtkWidget *windowToWidget(GtkWindow *window);
extern GtkGLArea *gobjectToGLArea(GObject *object);
extern GtkListBox *widgetToListBox(GtkWidget *widget);
extern GtkLabel *widgetToLabel(GtkWidget *widget);
extern GtkWidget *labelToWidget(GtkLabel *label);
extern GtkListBox *gobjectToListBox(GObject *object);
extern GtkMenuItem *gobjectToMenuItem(GObject *object);
extern GtkWidget *menuItemToWidget(GtkMenuItem *menuItem);
extern GtkGLArea *widgetToGLArea(GtkWidget *widget);
extern GtkWidget *listBoxToWidget(GtkListBox *listBox);
extern GtkContainer *listBoxToContainer(GtkListBox *listBox);
extern GtkWidget *listBoxRowToWidget(GtkListBoxRow *listBoxRow);
extern GtkContainer *listBoxRowToContainer(GtkListBoxRow *listBoxRow);
extern GtkWidget *toolButtonToWidget(GtkToolButton *toolButton);
extern GtkToolButton *gobjectToToolButton(GObject *object);
extern GtkBox *gobjectToBox(GObject *object);
extern GtkFileChooserDialog *widgetToFileChooserDialog(GtkWidget *widget);
extern GtkDialog *fileChooserDialogToDialog(GtkFileChooserDialog *fcdialog);
extern GtkWidget *fileChooserDialogToWidget(GtkFileChooserDialog *fcdialog);
extern GtkWidget *dialogToWidget(GtkDialog *dialog);
extern GtkFileChooser *fileChooserDialogToFileChooser(GtkFileChooserDialog *fcdialog);
extern GObject *labelToGObject(GtkLabel *label);
extern const gchar *gpointerToGChar(const gpointer data);
extern const gpointer gcharToGPointer(const gchar *data);
extern GtkWidget *imageToWidget(GtkImage *img);
extern GtkImage *widgetToImage(GtkWidget *widget);
extern GBytes *voidpToGbytes(void *data);
extern GtkMenuItem *widgetToMenuItem(GtkWidget *widget);
extern GtkContainer *menuBarToContainer(GtkMenuBar *menuBar);
extern GtkMenu *widgetToMenu(GtkWidget *widget);
extern GtkMenuShell *menuToMenuShell(GtkMenu *menu);
extern GtkMenuShell *menuBarToMenuShell(GtkMenuBar *menuBar);
extern GtkWidget *menuBarToWidget(GtkMenuBar *menuBar);
extern GtkMenuBar *widgetToMenuBar(GtkWidget *widget);
extern GtkWidget *menuToWidget(GtkMenu *menu);
extern GtkWidget *entryToWidget(GtkEntry *entry);
extern GtkEntry *widgetToEntry(GtkWidget *widget);
extern GtkEditable *entryToEditable(GtkEntry *entry);
extern GdkEventKey *eventToEventKey(GdkEvent *event);
extern GtkSwitch *widgetToSwitch(GtkWidget* widget);
extern GtkWidget *switchToWidget(GtkSwitch* gswitch);
extern GtkSpinButton *widgetToSpinButton(GtkWidget *widget);
extern GtkWidget *spinButtonToWidget(GtkSpinButton *spinButton);

extern void signalConnectButton(GtkButton *button, const char *signal, int id);
extern void sizeAllocateSignalConnectWidget(GtkWidget *widget, const char *signal, const char *name);
extern void signalConnectMenuItem(GtkMenuItem *menuItem, const char *signal, const char *name);
extern void eventSignalConnectWidget(GtkWidget *widget, const char *signal, const char *name);
extern void rowSelectedSignalConnectListBox(GtkListBox *listBox, const char *signal, const char *name);
extern void signalConnectToolButton(GtkToolButton *toolButton, const char *name);
extern void signalConnectWidget(GtkWidget *widget, const char *signal, const char *name);

extern GtkWidget *gohome_file_chooser_dialog_new(const gchar *title, GtkWindow *parent, GtkFileChooserAction action);

#endif