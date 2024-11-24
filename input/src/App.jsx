import React, { useState, useEffect } from 'react';
import './App.css';
import Modal from 'react-modal';

const PALETTE = {
  gray: '#dddddd',
  darkBlue: '#4365b4',
  purple: '#8275bc',
  blue: '#03a4f5',
  green: '#56b49a',
  yellow: '#e7e27e',
  orange: '#dea383',
  red: '#da6e99',
  pink: '#c95bae',
}

const FACE_OPTIONS = {
  Border: ['gray/gray'],
  Circle: ['orange/darkBlue', 'green/darkBlue'],
  Octogon: ['yellow/green', 'pink/darkBlue'],
  BeveledSquare: ['yellow/orange', 'b:darkBlue/pink'],
  Gate: ['orange/yellow', 'orange/blue'],
  Triangle: ['green/yellow', 'red/green', 'b:yellow/red'],
  Plane: ['yellow/red', 'purple/green'],
  Crown: ['orange/purple', 'pink/green', 'b:darkBlue/blue'],
  Error: ['purple/blue', 'purple/green', 's:green/darkBlue'],
  Tetris: ['blue/purple', 'b:green/red'],
  Potion: ['blue/red', 'yellow/darkBlue'],
  Bird: ['red/yellow', 'orange/green', 'b:blue/pink'],
  Arrows: ['red/blue', 'b:pink/yellow'],
  Shuriken: ['green/blue', 'purple/orange', 'red/blue'],
  Star: ['yellow/blue', 'pink/green'],
}

const App = () => {

  const [pieces, setPieces] = useState(JSON.parse(localStorage.getItem('geomtrix-pieces-v1')) ?? []);

  useEffect(() => {
    if (pieces.length) {
      const counts = {}
      pieces.forEach(piece => {
        ['north','south','east','west'].forEach(face => {
          const key = `${piece[face].kind}/${piece[face].fgColor}/${piece[face].bgColor}`;
          counts[key] = (counts[key] || 0) + 1;
        });
      });

      const notInOptions = {...counts};
      Object.keys(FACE_OPTIONS).forEach(kind => {
        FACE_OPTIONS[kind].forEach(colorCombination => {
          const key = `${kind}/${colorCombination.replace('b:', '').replace('s:', '')}`;
          delete notInOptions[key];
          if (!counts[key]) {
            console.warn('unused combination:', kind, colorCombination, key, counts);
          }
        });
      })

      Object.entries(notInOptions).forEach(([key, count]) => {
        console.warn('impossible combination:', key, count);
      });
    }
  }, []);

  useEffect(() => {
    localStorage.setItem('geomtrix-pieces-v1', JSON.stringify(pieces));
    if (pieces.length % 20 === 0) {
      localStorage.setItem(`geometrix-${pieces.length}`, JSON.stringify(pieces));
    }
  }, [pieces]);

  const onPieceFinished = (piece) => {
    setPieces([...pieces, piece]);
  }

  const onPieceDropped = (piece) => {
    setPieces(pieces.filter(p => p !== piece));
  }

  return (
    <div>
      <h3>Catálogo de Peças</h3>
      <hr />
      { pieces.length < 400 && <PieceBuilder currentPiece={pieces.length + 1} onPieceFinished={onPieceFinished} />}
      { pieces.length == 400 && (
        <div className="catalog-output">
          <pre>
            {catalogOutput(pieces)}
          </pre>
        </div>
      )}
      <hr />
      <div className="finished-pieces">
        {[...pieces].reverse().map((piece, index) => (
          <Piece key={index} number={pieces.length - index} onRemove={index === 0 ? () => onPieceDropped(piece) : null}>
            {partFromSpec(piece.north, 'north')}
            {partFromSpec(piece.south, 'south')}
            {partFromSpec(piece.east, 'east')}
            {partFromSpec(piece.west, 'west')}
          </Piece>
        ))}
      </div>
    </div>
  );
};

const PieceBuilder = ({ currentPiece, onPieceFinished }) => {
  const [north, setNorth] = useState(currentPiece >= 325 ? { kind: 'Border', fgColor: 'gray', bgColor: 'gray' } : null);
  const [east, setEast] = useState(currentPiece >= 397 ? { kind: 'Border', fgColor: 'gray', bgColor: 'gray' } : null);
  const [south, setSouth] = useState(null);
  const [west, setWest] = useState(null);
  const [activeFace, setActiveFace] = useState(null);

  const handleMouseUp = (event) => {
    const x = event.nativeEvent.offsetX;
    const y = event.nativeEvent.offsetY;

    if (triangleContains(0, 0, 80, 80, 160, 0, x, y)) {
      setActiveFace('north');
    } else if (triangleContains(0, 160, 80, 80, 160, 160, x, y)) {
      setActiveFace('south');
    } else if (triangleContains(160, 0, 80, 80, 160, 160, x, y)) {
      setActiveFace('east');
    } else if (triangleContains(0, 0, 80, 80, 0, 160, x, y)) {
      setActiveFace('west');
    }
  };

  const updateActiveFace = ({ kind, fgColor, bgColor }) => {
    const setter = {
      north: setNorth,
      south: setSouth,
      east: setEast,
      west: setWest,
    }[activeFace];

    setter({ kind, fgColor, bgColor });

    const ordering = ['north', 'east', 'south', 'west', 'north'];
    const nextFace = ordering.slice(ordering.indexOf(activeFace) + 1).find(face => face !== activeFace && !({ north, south, east, west })[face]);
    if (nextFace) { setActiveFace(nextFace); }
    else { setActiveFace(null); }
  };

  useEffect(() => {
    if (north && south && east && west) {
      onPieceFinished({ north, south, east, west });
      setNorth(currentPiece + 1 >= 325 ? { kind: 'Border', fgColor: 'gray', bgColor: 'gray' } : null);
      setEast(currentPiece + 1 >= 397 ? { kind: 'Border', fgColor: 'gray', bgColor: 'gray' } : null);
      setSouth(null);
      setWest(null);
    }
  });

  const faceOptions = [];
  Object.entries(FACE_OPTIONS).forEach(([kind, colorCombinations]) => {
    colorCombinations.forEach(colorCombination => {
      const requiresBorder = currentPiece >= 325 && activeFace !== 'south' || currentPiece >= 397;
      const isBorder = colorCombination.startsWith('b:') || colorCombination.startsWith('s:');
      if (requiresBorder !== isBorder && !colorCombination.startsWith('s:')) return;
      const [fgColor, bgColor] = colorCombination.replace('b:', '').replace('s:', '').split('/');
      faceOptions.push({ kind, fgColor, bgColor });
    });
  });
  faceOptions.sort((a, b) => {
    return Object.keys(PALETTE).indexOf(a.bgColor) - Object.keys(PALETTE).indexOf(b.bgColor);
  });

  return (
    <div className="piece-builder" style={{ position: "relative", display: "inline-block" }}>
      <Piece number={currentPiece}>
        {north && partFromSpec(north, 'north')}
        {south && partFromSpec(south, 'south')}
        {east && partFromSpec(east, 'east')}
        {west && partFromSpec(west, 'west')}
      </Piece>
      <div className="piece-clicker" style={{ cursor: 'pointer', position: "absolute", top: 0, left: 0, width: '100%', height: '100%' }} onMouseUp={handleMouseUp}></div>
      <Modal ariaHideApp={false} isOpen={!!activeFace} shouldCloseOnOverlayClick onRequestClose={() => setActiveFace(null)}>
        <div className="picker-table" data-orientation={activeFace}>
          {faceOptions.map(({ kind, fgColor, bgColor }) => (
            <div className="picker-cell" key={`${kind}-${fgColor}-${bgColor}`} onClick={() => updateActiveFace({ kind, fgColor, bgColor })}>
              {partFromSpec({ kind, fgColor, bgColor }, activeFace ?? 'north')}
            </div>
          ))}
        </div>
      </Modal>
    </div>
  )
}

const partFromSpec = (spec, orientation) => {
  const { kind, fgColor, bgColor } = spec;

  const components = {
    Border,
    Circle,
    BeveledSquare,
    Triangle,
    Plane,
    Crown,
    Error,
    Tetris,
    Potion,
    Octogon,
    Bird,
    Arrows,
    Shuriken,
    Star,
    Gate,
  }
  const Component = components[kind];

  return <Component bgColor={PALETTE[bgColor]} fgColor={PALETTE[fgColor]} orientation={orientation} />
}

export default App;

const Piece = ({ children, number, onRemove }) => {
  if (number) {
    return (
      <div className="numbered-piece">
        <div className="number">
          #{number}

          { onRemove ? (
            <>{' '}<a href="javascript:void(0)" style={{ cursor: 'pointer', textDecoration: 'none' }} onClick={onRemove}>🗑</a></>
          ) : null }
        </div>
        <Piece>{children}</Piece>
      </div>
    )
  }

  return (
    <div className="piece">
      {children}
      <div style={{ position: 'absolute', top: 0, left: 0, width: '100%', height: '100%' }}>
        <svg width="160" height="160" xmlns="http://www.w3.org/2000/svg">
          <line x1="0" y1="0" x2="160" y2="160" stroke="black" strokeWidth="1" />
          <line x1="160" y1="0" x2="0" y2="160" stroke="black" strokeWidth="1" />
        </svg>
      </div>
    </div>
  )
}

const PieceCore = ({ bgColor, fgColor, orientation, children }) => {
  return (
    <div className={`piece-part orientation-${orientation}`}>
      <svg width="160" height="80" xmlns="http://www.w3.org/2000/svg">
        <polygon points="0,0 80,80 160,0" fill={bgColor} />
        {children({ bgColor, fgColor })}
      </svg>
    </div>
  )
}

const Border = (props) => <PieceCore {...props}>{() => null}</PieceCore>

const Circle = (props) => {
  return (
    <PieceCore {...props}>
      {({ fgColor }) => <path d="M 126,0 A 46,46 0 0 1 34,0 M 46,0 A 34,34 0 0 0 114,0" stroke="black" strokeWidth="1" fill={fgColor} />}
    </PieceCore>
  )
}

const BeveledSquare = (props) => {
  return (
    <PieceCore {...props}>
      {({ fgColor }) => <path d="M 40,0 L 40,20 A 20,20 0 0 1 60,40 L 100,40 A 20,20 0 0 1 120,20 L 120,0" stroke="black" strokeWidth="1" fill={fgColor} />}
    </PieceCore>
  )
}

const Triangle = (props) => {
  return (
    <PieceCore {...props}>
      {({ fgColor }) => <path d="M 20,0 L 80,60 L 140,0 M 116,0 L 80,36 L 44,0" stroke="black" strokeWidth="1" fill={fgColor} />}
    </PieceCore>
  )
}

const Plane = (props) => {
  return (
    <PieceCore {...props}>
      {({ fgColor }) => (
        <path d="M 24,0 L 68,12 L 80,56 L 92,12 L 136,0" stroke="black" strokeWidth="1" fill={fgColor} />
      )}
    </PieceCore>
  )
}

const Crown = (props) => {
  return (
    <PieceCore {...props}>
      {({ fgColor }) => (
        <path d="M 56,0 L 32,24 L 68,20 L 80,56 L 92,20 L 128,24 L 104,0" stroke="black" strokeWidth="1" fill={fgColor} />
      )}
    </PieceCore>
  )
}

const Gate = (props) => {
  return (
    <PieceCore {...props}>
      {({ fgColor }) => (
        <path d="M 40,0 L 40,20 A 20,20 0 0 1 60,40 L 100,40 A 20,20 0 0 1 120,20 L 120,0 M 104,0 A 24,24 0 0 1 56,0" stroke="black" strokeWidth="1" fill={fgColor} />
      )}
    </PieceCore>
  )
}

const Octogon = (props) => {
  return (
    <PieceCore {...props}>
      {({ fgColor }) => (
        <path d="M 37,0 L 37,18 L 62,43 L 98,43 L 123,18 L 123,0 M 110,0 A 30,30 0 0 1 50,0" stroke="black" strokeWidth="1" fill={fgColor} />
      )}
    </PieceCore>
  )
}

const Arrows = (props) => {
  return (
    <PieceCore {...props}>
      {({ fgColor }) => (
        <path d="M 27,0 L 43,28 L 43,18 L 53,18 L 53,28 L 63,28 L 63,38 L 53,38 L 80,53 L 107,38 L 97,38 L 97,28 L 107,28 L 107,18 L117,18 L 117,28 L 133,0" stroke="black" strokeWidth="1" fill={fgColor} />
      )}
    </PieceCore>
  )
}

const Shuriken = (props) => {
  return (
    <PieceCore {...props}>
      {({ fgColor }) => (
        <path d="M 27,0 L 43,28 L 43,18 L 53,18 L 53,28 L 63,28 L 63,38 L 53,38 L 80,53 L 107,38 L 97,38 L 97,28 L 107,28 L 107,18 L117,18 L 117,28 L 133,0 M 100,0 A 20,20 0 0 1 60,0" stroke="black" strokeWidth="1" fill={fgColor} />
      )}
    </PieceCore>
  )
}

const Tetris = (props) => {
  return (
    <PieceCore {...props}>
      {({ fgColor }) => (
        <path d="M 50,0 L 36,14 L 66,44 L 80,30 L 94,44 L 124,14 L 110,0" stroke="black" strokeWidth="1" fill={fgColor} />
      )}
    </PieceCore>
  )
}

const Error = (props) => {
  return (
    <PieceCore {...props}>
      {({ fgColor }) => (
        <path d="M 36,0 A 44,44 0 0 0 124,0 M 97,0 L 108,11 L 91,28 L 80,17 L 69,28 L 52,11 L 63,0" stroke="black" strokeWidth="1" fill={fgColor} />
      )}
    </PieceCore>
  )
}

const Potion = (props) => {
  return (
    <PieceCore {...props}>
      {({ fgColor }) => (
        <path d="M 32,0 L 68,19 L 56,47 L 104,47 L 92,19 L 128,0" stroke="black" strokeWidth="1" fill={fgColor} />
      )}
    </PieceCore>
  )
}

const Bird = (props) => {
  return (
    <PieceCore {...props}>
      {({ fgColor }) => (
        <path d="M 12,0 L 50,22 L 68,10 L 56,34 L 80,68 L 104,34 L 92,10 L 110,22 L 148,0" stroke="black" strokeWidth="1" fill={fgColor} />
      )}
    </PieceCore>
  )
}

const Star = (props) => {
  return (
    <PieceCore {...props}>
      {({ fgColor }) => (
        <path d="M 28,0 L 58,8 L 44,34 L 72,21 L 80,52 L 88,21 L 116,34 L 102,8 L 132,0" stroke="black" strokeWidth="1" fill={fgColor} />
      )}
    </PieceCore>
  )
}

function triangleContains(ax, ay, bx, by, cx, cy, x, y) {

  let det = (bx - ax) * (cy - ay) - (by - ay) * (cx - ax)

  return  det * ((bx - ax) * (y - ay) - (by - ay) * (x - ax)) >= 0 &&
          det * ((cx - bx) * (y - by) - (cy - by) * (x - bx)) >= 0 &&
          det * ((ax - cx) * (y - cy) - (ay - cy) * (x - cx)) >= 0

}

const catalogOutput = (pieces) => {
  const possibleFaces = {};

  pieces.forEach(piece => {
    ['north','south','east','west'].forEach(face => {
      const key = `${piece[face].kind}/${piece[face].fgColor}/${piece[face].bgColor}`;
      possibleFaces[key] = true;
    });
  });

  const sortedKeys = Object.keys(possibleFaces).sort((a, b) => a.localeCompare(b));

  Object.keys(possibleFaces).forEach(key => {
    possibleFaces[key] = sortedKeys.indexOf(key);
  });

  return '[\n' + pieces.map(piece => (
    '  ' + ['north','south','east','west'].map(face => possibleFaces[`${piece[face].kind}/${piece[face].fgColor}/${piece[face].bgColor}`]).join(', ')
  )).join(',\n') + '\n]';
}
