package jsonmodel

import (
	"testing"
)

// TestExtendPath vérifie la création correcte de la liste chaînée du chemin
func TestExtendPath(t *testing.T) {
	root := &Path{Name: "root", Index: -1}
	child := ExtendPath(root, "sub")

	if child.Parent != root {
		t.Errorf("Le parent devrait être root")
	}
	if child.Name != "sub" {
		t.Errorf("Le nom attendu était 'sub', obtenu %s", child.Name)
	}
	if child.Index != -1 {
		t.Errorf("L'index par défaut pour une propriété doit être -1")
	}
}

// TestSelectPath valide le retour conditionnel pour le reporting
func TestSelectPath(t *testing.T) {
	p := &Path{Name: "test"}
	
	if SelectPath(p, true) != p {
		t.Error("SelectPath(true) devrait retourner le chemin")
	}
	if SelectPath(p, false) != nil {
		t.Error("SelectPath(false) devrait retourner nil")
	}
}

// TestLen vérifie la réflexion sur différentes structures
func TestLen(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected int
	}{
		{map[string]interface{}{"a": 1, "b": 2}, 2},
		{[]int{1, 2, 3, 4}, 4},
		{"hello", 5},
		{123, 0}, // Type non supporté par Len()
		{nil, 0},
	}

	for _, tt := range tests {
		result := Len(tt.input)
		if result != tt.expected {
			t.Errorf("Len(%v): attendu %d, obtenu %d", tt.input, tt.expected, result)
		}
	}
}

// TestObjectHasPropVal vérifie l'extraction sécurisée des propriétés
func TestObjectHasPropVal(t *testing.T) {
	obj := map[string]interface{}{"name": "Hobbes"}
	var dst interface{}

	// Cas : Propriété existante
	if !ObjectHasPropVal(obj, "name", &dst) {
		t.Error("Devrait trouver la propriété 'name'")
	}
	if dst != "Hobbes" {
		t.Errorf("Valeur incorrecte: attendu Hobbes, obtenu %v", dst)
	}

	// Cas : Propriété manquante
	if ObjectHasPropVal(obj, "age", &dst) {
		t.Error("Ne devrait pas trouver la propriété 'age'")
	}

	// Cas : Objet de type invalide (pas une map)
	if ObjectHasPropVal([]string{"pas", "une", "map"}, "name", &dst) {
		t.Error("Devrait retourner false pour un type non-map")
	}
}