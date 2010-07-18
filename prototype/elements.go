package main


type Element interface {
        GetName() string
        GetDamageMod(Element) float // Returns the damage mod for self->target, eg for Electric.GetDamageMod( Earth ) -> 0
        SetDamageMod(Element, float)                                                                                      
}                                                                                                                         

type ElementList interface {
        GetDamageMod(Element) float // Returns the damage mod for self->target, eg for Electric|Water GetDamageMod( Fire ) -> 1*2
        AddElements(...Element)                                                                                                  
        RemoveElements(...Element)                                                                                               
        Iter() <-chan Element                                                                                                    
}                 


//=============================================================================================================================================
//                                                              element                                                                        
//=============================================================================================================================================
type element struct {                                                                                                                          
        name      string                                                                                                                       
        damageMod map[Element]float                                                                                                            
}                                                                                                                                              

func TestElement(name string) Element {
        newElement := new(element)    
        newElement.name = name        
        newElement.damageMod = make(map[Element]float)
        return newElement                             
}                                                     

func (e *element) GetName() string {
        return e.name               
}                                   

func (e *element) GetDamageMod(target Element) float {
        if mod, ok := e.damageMod[target]; ok {       
                return mod                            
        }                                             
        return 1.0                                    
}                                                     

func (e *element) SetDamageMod(target Element, mod float) {
        e.damageMod[target] = mod                          
}                                                          

//=============================================================================================================================================
//                                                              elementList                                                                    
//=============================================================================================================================================
type elementList struct {                                                                                                                      
        elements  map[Element]bool                                                                                                             
        damageMod map[Element]float                                                                                                            
}                                                                                                                                              

func TestElementList(elements ...Element) ElementList {
        newList := new(elementList)                   
        newList.elements = make(map[Element]bool)     
        newList.damageMod = make(map[Element]float)   
        newList.AddElements(elements)                 
        return newList                                
}                                                     

func (l *elementList) iterate(c chan<- Element) {
        for e, hasE := range l.elements {        
                if hasE {                        
                        c <- e                   
                }                                
        }                                        
        close(c)                                 
}                                                

func (l *elementList) clearMods() {
        l.damageMod = make(map[Element]float)
}                                            

func (l *elementList) Iter() <-chan Element {
        c := make(chan Element)              
        go l.iterate(c)                      
        return c                             
}                                            


func (l *elementList) AddElements(elements ...Element) {
        for _, e := range elements {                    
                l.elements[e] = true                    
        }                                               
        l.clearMods()                                   
}                                                       

func (l *elementList) RemoveElements(elements ...Element) {
        for _, e := range elements {                       
                l.elements[e] = false                      
        }                                                  
        l.clearMods()                                      
}                                                          


func (l *elementList) GetDamageMod(target Element) float {
        if mod, ok := l.damageMod[target]; ok {           
                return mod                                
        }                                                 
        mod := 1.0                                        
        for e := range l.Iter() {                         
                mod *= e.GetDamageMod(target)             
        }                                                 
        return mod                                        
}                                                         

