package engine

import (
	"fmt"
	"strings"

	"github.com/seaweedfs/cockroachdb-parser/pkg/sql/parser"
	"github.com/seaweedfs/cockroachdb-parser/pkg/sql/sem/tree"
)

// CockroachSQLParser wraps CockroachDB's PostgreSQL-compatible SQL parser for use in Hanzo
type CockroachSQLParser struct{}

// NewCockroachSQLParser creates a new instance of the CockroachDB SQL parser wrapper
func NewCockroachSQLParser() *CockroachSQLParser {
	return &CockroachSQLParser{}
}

// ParseSQL parses a SQL statement using CockroachDB's parser
func (p *CockroachSQLParser) ParseSQL(sql string) (Statement, error) {
	// Parse using CockroachDB's parser
	stmts, err := parser.Parse(sql)
	if err != nil {
		return nil, fmt.Errorf("CockroachDB parser error: %v", err)
	}

	if len(stmts) != 1 {
		return nil, fmt.Errorf("expected exactly one statement, got %d", len(stmts))
	}

	stmt := stmts[0].AST

	// Convert CockroachDB AST to Hanzo AST format
	switch s := stmt.(type) {
	case *tree.Select:
		return p.convertSelectStatement(s)
	default:
		return nil, fmt.Errorf("unsupported statement type: %T", s)
	}
}

// convertSelectStatement converts CockroachDB's Select AST to Hanzo format
func (p *CockroachSQLParser) convertSelectStatement(crdbSelect *tree.Select) (*SelectStatement, error) {
	selectClause, ok := crdbSelect.Select.(*tree.SelectClause)
	if !ok {
		return nil, fmt.Errorf("expected SelectClause, got %T", crdbSelect.Select)
	}

	hanzoSelect := &SelectStatement{
		SelectExprs: make([]SelectExpr, 0, len(selectClause.Exprs)),
		From:        []TableExpr{},
	}

	// Convert SELECT expressions
	for _, expr := range selectClause.Exprs {
		hanzoExpr, err := p.convertSelectExpr(expr)
		if err != nil {
			return nil, fmt.Errorf("failed to convert select expression: %v", err)
		}
		hanzoSelect.SelectExprs = append(hanzoSelect.SelectExprs, hanzoExpr)
	}

	// Convert FROM clause
	if len(selectClause.From.Tables) > 0 {
		for _, fromExpr := range selectClause.From.Tables {
			hanzoTableExpr, err := p.convertFromExpr(fromExpr)
			if err != nil {
				return nil, fmt.Errorf("failed to convert FROM clause: %v", err)
			}
			hanzoSelect.From = append(hanzoSelect.From, hanzoTableExpr)
		}
	}

	// Convert WHERE clause if present
	if selectClause.Where != nil {
		whereExpr, err := p.convertExpr(selectClause.Where.Expr)
		if err != nil {
			return nil, fmt.Errorf("failed to convert WHERE clause: %v", err)
		}
		hanzoSelect.Where = &WhereClause{
			Expr: whereExpr,
		}
	}

	// Convert LIMIT and OFFSET clauses if present
	if crdbSelect.Limit != nil {
		limitClause := &LimitClause{}

		// Convert LIMIT (Count)
		if crdbSelect.Limit.Count != nil {
			countExpr, err := p.convertExpr(crdbSelect.Limit.Count)
			if err != nil {
				return nil, fmt.Errorf("failed to convert LIMIT clause: %v", err)
			}
			limitClause.Rowcount = countExpr
		}

		// Convert OFFSET
		if crdbSelect.Limit.Offset != nil {
			offsetExpr, err := p.convertExpr(crdbSelect.Limit.Offset)
			if err != nil {
				return nil, fmt.Errorf("failed to convert OFFSET clause: %v", err)
			}
			limitClause.Offset = offsetExpr
		}

		hanzoSelect.Limit = limitClause
	}

	return hanzoSelect, nil
}

// convertSelectExpr converts CockroachDB SelectExpr to Hanzo format
func (p *CockroachSQLParser) convertSelectExpr(expr tree.SelectExpr) (SelectExpr, error) {
	// Handle star expressions (SELECT *)
	if _, isStar := expr.Expr.(tree.UnqualifiedStar); isStar {
		return &StarExpr{}, nil
	}

	// CockroachDB's SelectExpr is a struct, not an interface, so handle it directly
	hanzoExpr := &AliasedExpr{}

	// Convert the main expression
	convertedExpr, err := p.convertExpr(expr.Expr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert expression: %v", err)
	}
	hanzoExpr.Expr = convertedExpr

	// Convert alias if present
	if expr.As != "" {
		hanzoExpr.As = aliasValue(expr.As)
	}

	return hanzoExpr, nil
}

// convertExpr converts CockroachDB expressions to Hanzo format
func (p *CockroachSQLParser) convertExpr(expr tree.Expr) (ExprNode, error) {
	switch e := expr.(type) {
	case *tree.FuncExpr:
		// Function call
		hanzoFunc := &FuncExpr{
			Name:  stringValue(strings.ToUpper(e.Func.String())), // Convert to uppercase for consistency
			Exprs: make([]SelectExpr, 0, len(e.Exprs)),
		}

		// Convert function arguments
		for _, arg := range e.Exprs {
			// Special case: Handle star expressions in function calls like COUNT(*)
			if _, isStar := arg.(tree.UnqualifiedStar); isStar {
				hanzoFunc.Exprs = append(hanzoFunc.Exprs, &StarExpr{})
			} else {
				convertedArg, err := p.convertExpr(arg)
				if err != nil {
					return nil, fmt.Errorf("failed to convert function argument: %v", err)
				}
				hanzoFunc.Exprs = append(hanzoFunc.Exprs, &AliasedExpr{Expr: convertedArg})
			}
		}

		return hanzoFunc, nil

	case *tree.BinaryExpr:
		// Arithmetic/binary operations (including string concatenation ||)
		hanzoArith := &ArithmeticExpr{
			Operator: e.Operator.String(),
		}

		// Convert left operand
		left, err := p.convertExpr(e.Left)
		if err != nil {
			return nil, fmt.Errorf("failed to convert left operand: %v", err)
		}
		hanzoArith.Left = left

		// Convert right operand
		right, err := p.convertExpr(e.Right)
		if err != nil {
			return nil, fmt.Errorf("failed to convert right operand: %v", err)
		}
		hanzoArith.Right = right

		return hanzoArith, nil

	case *tree.ComparisonExpr:
		// Comparison operations (=, >, <, >=, <=, !=, etc.) used in WHERE clauses
		hanzoComp := &ComparisonExpr{
			Operator: e.Operator.String(),
		}

		// Convert left operand
		left, err := p.convertExpr(e.Left)
		if err != nil {
			return nil, fmt.Errorf("failed to convert comparison left operand: %v", err)
		}
		hanzoComp.Left = left

		// Convert right operand
		right, err := p.convertExpr(e.Right)
		if err != nil {
			return nil, fmt.Errorf("failed to convert comparison right operand: %v", err)
		}
		hanzoComp.Right = right

		return hanzoComp, nil

	case *tree.StrVal:
		// String literal
		return &SQLVal{
			Type: StrVal,
			Val:  []byte(string(e.RawString())),
		}, nil

	case *tree.NumVal:
		// Numeric literal
		valStr := e.String()
		if strings.Contains(valStr, ".") {
			return &SQLVal{
				Type: FloatVal,
				Val:  []byte(valStr),
			}, nil
		} else {
			return &SQLVal{
				Type: IntVal,
				Val:  []byte(valStr),
			}, nil
		}

	case *tree.UnresolvedName:
		// Column name
		return &ColName{
			Name: stringValue(e.String()),
		}, nil

	case *tree.AndExpr:
		// AND expression
		left, err := p.convertExpr(e.Left)
		if err != nil {
			return nil, fmt.Errorf("failed to convert AND left operand: %v", err)
		}
		right, err := p.convertExpr(e.Right)
		if err != nil {
			return nil, fmt.Errorf("failed to convert AND right operand: %v", err)
		}
		return &AndExpr{
			Left:  left,
			Right: right,
		}, nil

	case *tree.OrExpr:
		// OR expression
		left, err := p.convertExpr(e.Left)
		if err != nil {
			return nil, fmt.Errorf("failed to convert OR left operand: %v", err)
		}
		right, err := p.convertExpr(e.Right)
		if err != nil {
			return nil, fmt.Errorf("failed to convert OR right operand: %v", err)
		}
		return &OrExpr{
			Left:  left,
			Right: right,
		}, nil

	case *tree.Tuple:
		// Tuple expression for IN clauses: (value1, value2, value3)
		tupleValues := make(ValTuple, 0, len(e.Exprs))
		for _, tupleExpr := range e.Exprs {
			convertedExpr, err := p.convertExpr(tupleExpr)
			if err != nil {
				return nil, fmt.Errorf("failed to convert tuple element: %v", err)
			}
			tupleValues = append(tupleValues, convertedExpr)
		}
		return tupleValues, nil

	case *tree.CastExpr:
		// Handle INTERVAL expressions: INTERVAL '1 hour'
		// CockroachDB represents these as cast expressions
		if p.isIntervalCast(e) {
			// Extract the string value being cast to interval
			if strVal, ok := e.Expr.(*tree.StrVal); ok {
				return &IntervalExpr{
					Value: string(strVal.RawString()),
				}, nil
			}
			return nil, fmt.Errorf("invalid INTERVAL expression: expected string literal")
		}
		// For non-interval casts, just convert the inner expression
		return p.convertExpr(e.Expr)

	case *tree.RangeCond:
		// Handle BETWEEN expressions: column BETWEEN value1 AND value2
		hanzoBetween := &BetweenExpr{
			Not: e.Not, // Handle NOT BETWEEN
		}

		// Convert the left operand (the expression being tested)
		left, err := p.convertExpr(e.Left)
		if err != nil {
			return nil, fmt.Errorf("failed to convert BETWEEN left operand: %v", err)
		}
		hanzoBetween.Left = left

		// Convert the FROM operand (lower bound)
		from, err := p.convertExpr(e.From)
		if err != nil {
			return nil, fmt.Errorf("failed to convert BETWEEN from operand: %v", err)
		}
		hanzoBetween.From = from

		// Convert the TO operand (upper bound)
		to, err := p.convertExpr(e.To)
		if err != nil {
			return nil, fmt.Errorf("failed to convert BETWEEN to operand: %v", err)
		}
		hanzoBetween.To = to

		return hanzoBetween, nil

	case *tree.IsNullExpr:
		// Handle IS NULL expressions: column IS NULL
		expr, err := p.convertExpr(e.Expr)
		if err != nil {
			return nil, fmt.Errorf("failed to convert IS NULL expression: %v", err)
		}

		return &IsNullExpr{
			Expr: expr,
		}, nil

	case *tree.IsNotNullExpr:
		// Handle IS NOT NULL expressions: column IS NOT NULL
		expr, err := p.convertExpr(e.Expr)
		if err != nil {
			return nil, fmt.Errorf("failed to convert IS NOT NULL expression: %v", err)
		}

		return &IsNotNullExpr{
			Expr: expr,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported expression type: %T", e)
	}
}

// convertFromExpr converts CockroachDB FROM expressions to Hanzo format
func (p *CockroachSQLParser) convertFromExpr(expr tree.TableExpr) (TableExpr, error) {
	switch e := expr.(type) {
	case *tree.TableName:
		// Simple table name
		tableName := TableName{
			Name: stringValue(e.Table()),
		}

		// Extract database qualifier if present

		if e.Schema() != "" {
			tableName.Qualifier = stringValue(e.Schema())
		}

		return &AliasedTableExpr{
			Expr: tableName,
		}, nil

	case *tree.AliasedTableExpr:
		// Handle aliased table expressions (which is what CockroachDB uses for qualified names)
		if tableName, ok := e.Expr.(*tree.TableName); ok {
			hanzoTableName := TableName{
				Name: stringValue(tableName.Table()),
			}

			// Extract database qualifier if present
			if tableName.Schema() != "" {
				hanzoTableName.Qualifier = stringValue(tableName.Schema())
			}

			return &AliasedTableExpr{
				Expr: hanzoTableName,
			}, nil
		}

		return nil, fmt.Errorf("unsupported expression in AliasedTableExpr: %T", e.Expr)

	default:
		return nil, fmt.Errorf("unsupported table expression type: %T", e)
	}
}

// isIntervalCast checks if a CastExpr is casting to an INTERVAL type
func (p *CockroachSQLParser) isIntervalCast(castExpr *tree.CastExpr) bool {
	// Check if the target type is an interval type
	// CockroachDB represents interval types in the Type field
	// We need to check if it's an interval type by examining the type structure
	if castExpr.Type != nil {
		// Try to detect interval type by examining the AST structure
		// Since we can't easily access the type string, we'll be more conservative
		// and assume any cast expression on a string literal could be an interval
		if _, ok := castExpr.Expr.(*tree.StrVal); ok {
			// This is likely an INTERVAL expression since CockroachDB
			// represents INTERVAL '1 hour' as casting a string to interval type
			return true
		}
	}
	return false
}
